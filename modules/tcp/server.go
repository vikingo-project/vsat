package vstcp

import (
	"encoding/hex"
	"log"
	"net"
	"sync"
	"time"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
	API      *events.EventsAPI
	settings settings
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) handleConection(conn net.Conn) {

	clientIP := utils.ExtractIP(conn.RemoteAddr().String())
	session, _ := s.API.NewSession(models.SessionInfo{
		Description: "connect",
		ClientIP:    clientIP,
		LocalAddr:   conn.LocalAddr().String(),
	})
	if s.settings.Mode == "proxy" {
		raddr, err := net.ResolveTCPAddr("tcp", s.settings.ProxySettings.Destination)
		if err != nil {
			utils.PrintDebug("Failed to resolve remote address: %s", err)
			return
		}

		p := s.newProxy(session, conn, raddr)
		p.Start()
	} else {
		err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return
		}

		buf := make([]byte, 0xffff)
		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				conn.Close()
				return
			}
		}
		fields := make(map[string]interface{})
		conn.Write([]byte(s.settings.ResponseSettings.Response))
		conn.Close()

		if s.settings.LogRequest {
			fields["hex:request data"] = hex.EncodeToString(buf[:n])
		}
		if s.settings.LogResponse {
			fields["hex:response data"] = hex.EncodeToString([]byte(s.settings.ResponseSettings.Response))
		}
		if s.settings.LogRequest || s.settings.LogResponse {
			s.API.PushEvent(models.Event{Session: session, Data: fields})
		}
	}
}
