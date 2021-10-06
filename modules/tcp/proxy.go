package vstcp

import (
	"encoding/hex"
	"io"
	"net"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

type Proxy struct {
	EventsAPI    *events.EventsAPI
	settings     settings
	session      string
	raddr        *net.TCPAddr
	lconn, rconn io.ReadWriteCloser
	errsig       chan bool
}

func (s *Server) newProxy(session string, lconn net.Conn, raddr *net.TCPAddr) *Proxy {
	return &Proxy{
		settings:  s.settings,
		session:   session,
		EventsAPI: s.API,
		lconn:     lconn,
		raddr:     raddr,
		errsig:    make(chan bool),
	}
}

// Start - open connection to remote and start proxying data.
func (p *Proxy) Start() {
	defer p.lconn.Close()
	var err error
	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		utils.PrintDebug("Remote connection failed: %s", err)
		return
	}
	defer p.rconn.Close()
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)
	<-p.errsig
}

func (p *Proxy) pipe(src, dst io.ReadWriter) {
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			utils.PrintDebug("TCP proxy: failed to read TCP data %s\n", err.Error())
			p.errsig <- true
			return
		}
		b := buff[:n]
		if src == p.lconn && p.settings.LogRequest {
			p.EventsAPI.PushEvent(models.Event{Session: p.session, Data: map[string]interface{}{
				"hex:request data": hex.EncodeToString(b),
			}})
		} else if p.settings.LogResponse {
			p.EventsAPI.PushEvent(models.Event{Session: p.session, Data: map[string]interface{}{
				"hex:response data": hex.EncodeToString(b),
			}})
		}
		_, err = dst.Write(b)
		if err != nil {
			utils.PrintDebug("TCP proxy: failed to write to dst '%s'\n", err.Error())
			p.errsig <- true
			return
		}
	}
}
