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
	defer conn.Close()
	clientIP := utils.ExtractIP(conn.RemoteAddr().String())
	session, _ := s.API.NewSession(models.SessionInfo{
		Description: "connect",
		ClientIP:    clientIP,
		LocalAddr:   conn.LocalAddr().String(),
	})

	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			conn.Close()
			return
		}
	}

	s.API.PushEvent(models.Event{Session: session, Name: "open session", Data: map[string]interface{}{
		"hex:payload": hex.EncodeToString(buf[:n]),
	}})
	conn.Write([]byte("welcome to Vikingo Satellite"))
	conn.Close()
}
