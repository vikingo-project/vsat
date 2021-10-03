package vsdns

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

type Server struct {
	host      string
	port      int
	rTimeout  time.Duration
	wTimeout  time.Duration
	listeners []*dns.Server

	settings settings
	API      *events.EventsAPI
}

func (s *Server) Addr() string {
	return net.JoinHostPort(s.host, strconv.Itoa(s.port))
}

func (s *Server) Run() {
	handler := s.newHandler(s.API)

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", handler.DoTCP)

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", handler.DoUDP)

	tcpServer := &dns.Server{Addr: s.Addr(),
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout}

	udpServer := &dns.Server{Addr: s.Addr(),
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout}

	s.listeners = append(s.listeners, udpServer)
	s.listeners = append(s.listeners, tcpServer)

	go s.start(udpServer)
	go s.start(tcpServer)
}

func (s *Server) Stop() error {
	for _, l := range s.listeners {
		err := l.ShutdownContext(context.Background())
		if err != nil {
			return err
		}
		log.Println("killed")
	}
	return nil
}

func (s *Server) start(ds *dns.Server) {
	log.Printf("Start %s listener on %s", ds.Net, s.Addr())
	err := ds.ListenAndServe()
	if err != nil {
		log.Printf("Start %s listener on %s failed:%s", ds.Net, s.Addr(), err.Error())
	}
}

type ResolvError struct {
	qname, net  string
	nameservers []string
}

func (e ResolvError) Error() string {
	errmsg := fmt.Sprintf("%s resolv failed on %s (%s)", e.qname, strings.Join(e.nameservers, "; "), e.net)
	return errmsg
}

type Resolver struct {
	config *dns.ClientConfig
}

// Lookup will ask each nameserver in top-to-bottom fashion, starting a new request
// in every second, and return as early as possbile (have an answer).
// It returns an error if no request has succeeded.
func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, err error) {
	c := &dns.Client{
		Net:          net,
		ReadTimeout:  r.Timeout(),
		WriteTimeout: r.Timeout(),
	}
	/*
		if net == "udp" && settings.ResolvConfig.SetEDNS0 {
			req = req.SetEdns0(65535, true)
		}
	*/

	qname := req.Question[0].Name

	res := make(chan *dns.Msg, 1)
	var wg sync.WaitGroup
	L := func(nameserver string) {
		defer wg.Done()
		r, rtt, err := c.Exchange(req, nameserver)
		if err != nil {
			utils.PrintDebug("%s socket error on %s", qname, nameserver)
			utils.PrintDebug("error:%s", err.Error())
			return
		}
		// If SERVFAIL happen, should return immediately and try another upstream resolver.
		// However, other Error code like NXDOMAIN is an clear response stating
		// that it has been verified no such domain existas and ask other resolvers
		// would make no sense. See more about #20
		if r != nil && r.Rcode != dns.RcodeSuccess {
			utils.PrintDebug("%s failed to get an valid answer on %s", qname, nameserver)
			if r.Rcode == dns.RcodeServerFailure {
				return
			}
		} else {
			utils.PrintDebug("%s resolv on %s (%s) ttl: %d", UnFqdn(qname), nameserver, net, rtt)
		}
		select {
		case res <- r:
		default:
		}
	}

	ticker := time.NewTicker(time.Duration(500) * time.Millisecond)
	defer ticker.Stop()
	// Start lookup on each nameserver top-down, in every second
	for _, nameserver := range r.Nameservers() {
		log.Println("send recursive to", nameserver)
		wg.Add(1)
		go L(nameserver)
		// but exit early, if we have an answer
		select {
		case r := <-res:
			return r, nil
		case <-ticker.C:
			continue
		}
	}
	// wait for all the namservers to finish
	wg.Wait()
	select {
	case r := <-res:
		return r, nil
	default:
		return nil, ResolvError{qname, net, r.Nameservers()}
	}

}

// Namservers return the array of nameservers, with port number appended.
// '#' in the name is treated as port separator, as with dnsmasq.
func (r *Resolver) Nameservers() (ns []string) {
	for _, server := range r.config.Servers {
		if i := strings.IndexByte(server, '#'); i > 0 {
			server = net.JoinHostPort(server[:i], server[i+1:])
		} else {
			// add port if not exists
			parts := strings.Split(server, ":")
			if len(parts) == 1 {
				server = net.JoinHostPort(server, "53")
			}
		}
		ns = append(ns, server)
	}
	return
}

func (r *Resolver) Timeout() time.Duration {
	return time.Duration(5) * time.Second
}

type Question struct {
	qname  string
	qtype  string
	qclass string
}

func (q *Question) String() string {
	return q.qname + " " + q.qclass + " " + q.qtype
}

type GODNSHandler struct {
	resolver *Resolver
	server   *Server
	records  Records
	API      *events.EventsAPI
}

func (s *Server) newHandler(API *events.EventsAPI) *GODNSHandler {
	var resolver *Resolver
	clientConfig := &dns.ClientConfig{
		Servers: s.settings.Resolvers,
	}
	resolver = &Resolver{clientConfig}
	records := parseRecords(s.settings.Records)
	return &GODNSHandler{resolver: resolver, server: s, API: API, records: records}
}

func (h *GODNSHandler) do(Net string, w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	Q := Question{UnFqdn(q.Name), dns.TypeToString[q.Qtype], dns.ClassToString[q.Qclass]}

	var clientIP net.IP
	if Net == "tcp" {
		clientIP = w.RemoteAddr().(*net.TCPAddr).IP
	} else {
		clientIP = w.RemoteAddr().(*net.UDPAddr).IP
	}
	session, _ := h.API.NewSession(models.SessionInfo{
		Description: fmt.Sprintf("lookup %s [%s]", Q.qname, Q.qtype),
		ClientIP:    clientIP.String(),
		LocalAddr:   w.LocalAddr().String(),
	})

	h.API.PushEvent(models.Event{Session: session, Name: "lookup", Data: map[string]interface{}{"resource": Q.qname, "query type": Q.qtype}})
	log.Printf("%s lookup　%s", clientIP, Q.String())

	switch q.Qtype {
	case dns.TypeA:
		IPs := h.records.LookupA(Q.qname)
		if len(IPs) > 0 {
			m := new(dns.Msg)
			m.SetReply(req)
			rr_header := dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			}
			for _, ip := range IPs {
				a := &dns.A{Hdr: rr_header, A: net.ParseIP(ip)}
				m.Answer = append(m.Answer, a)
			}
			w.WriteMsg(m)
		}
	case dns.TypeAAAA:
		IPs := h.records.LookupAAAA(Q.qname)
		if len(IPs) > 0 {
			m := new(dns.Msg)
			m.SetReply(req)
			rr_header := dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    3600,
			}
			for _, ip := range IPs {
				aaaa := &dns.AAAA{Hdr: rr_header, AAAA: net.ParseIP(ip)}
				m.Answer = append(m.Answer, aaaa)
			}
			w.WriteMsg(m)
		}
	case dns.TypeCNAME:
		target := h.records.LookupCNAME(Q.qname)
		if len(target) > 0 {
			m := new(dns.Msg)
			m.SetReply(req)
			rr_header := dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    3600,
			}
			m.Answer = append(m.Answer, &dns.CNAME{Hdr: rr_header, Target: dns.Fqdn(target)})
			w.WriteMsg(m)
		}
	case dns.TypeNS:
		nameservers := h.records.LookupNS(Q.qname)
		if len(nameservers) > 0 {
			m := new(dns.Msg)
			m.SetReply(req)
			header := dns.RR_Header{
				Name:   dns.Fqdn(Q.qname),
				Rrtype: dns.TypeCNAME,
				Class:  dns.TypeNS,
				Ttl:    300,
			}
			for _, ns := range nameservers {
				m.Ns = append(m.Ns, &dns.NS{Hdr: header, Ns: ns})
			}
			w.WriteMsg(m)
		}
	case dns.TypeMX, dns.TypeTXT:
		// todo
	}

	if h.server.settings.Recursive {
		log.Println("send recursive query")
		mesg, err := h.resolver.Lookup(Net, req)
		if err != nil {
			utils.PrintDebug("Resolve query error %s", err)
			dns.HandleFailed(w, req)
			return
		}
		w.WriteMsg(mesg)
	} else {
		dns.HandleFailed(w, req)
	}
}

func (h *GODNSHandler) DoTCP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("tcp", w, req)
}

func (h *GODNSHandler) DoUDP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("udp", w, req)
}

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}

type KeyNotFound struct {
	key string
}

func (e KeyNotFound) Error() string {
	return e.key + " " + "not found"
}

type Mesg struct {
	Msg    *dns.Msg
	Expire time.Time
}
