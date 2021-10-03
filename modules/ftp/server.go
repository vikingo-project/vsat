// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vsftp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"
	"strings"
	"sync"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

type Storage struct {
	FilesTree sync.Map
}

func (s *Storage) realPath(path string) string {
	paths := strings.Split(path, "/")
	return filepath.Join(append([]string{"/"}, paths...)...)
}

type Server struct {
	feats        string
	PublicIP     string
	PassivePorts string
	quit         chan interface{}
	listener     net.Listener
	vfs          VirtualFS
	wg           sync.WaitGroup
	settings     settings
	EventsAPI    *events.EventsAPI
}

func (s *Server) wrapConn(tcpConn net.Conn) *Conn {
	c := new(Conn)
	c.curDir = "/"
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.server = s
	c.vfs = s.vfs.Clone()
	c.EventsAPI = s.EventsAPI
	return c
}

func (s *Server) handleConection(conn *Conn) {
	session, _ := s.EventsAPI.NewSession(models.SessionInfo{LocalAddr: conn.conn.LocalAddr().String(), ClientIP: utils.ExtractIP(conn.conn.RemoteAddr().String())})
	conn.session = session
	conn.writeMessage(220, "Welcome!")
	for {
		line, err := conn.controlReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(conn.session, fmt.Sprint("read error:", err))
			}
			break
		}
		conn.receiveLine(line)
		if conn.closed {
			break
		}
	}
	conn.Close()
	log.Println(conn.session, "Connection Terminated")
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
				s.handleConection(s.wrapConn(conn))
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
}
