package tunnels

import (
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/vikingo-project/vsat/shared"
	"github.com/vikingo-project/vsat/utils"
	"github.com/vmihailenco/msgpack"
)

type Con struct {
	net.Conn
	rwm *sync.RWMutex
}

func (c *Con) ReadPacket() (action uint8, data []byte, err error) {
	if err = binary.Read(c, binary.LittleEndian, &action); err != nil {
		utils.PrintDebug("failed to read packet '%s'\n", err)
		return
	}

	utils.PrintDebug("ReadPacket# action %d", action)
	var length uint16
	if err = binary.Read(c, binary.LittleEndian, &length); err != nil {
		utils.PrintDebug("Read failed binary'%s'\n", err)
		return
	}
	data = make([]byte, length)
	io.ReadFull(c, data)
	return
}

func (c *Con) WritePacket(t uint8, msg interface{}) error {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	utils.PrintDebug("writepacket %d, %v", t, msg)
	var (
		b   []byte
		err error
	)

	b, err = msgpack.Marshal(msg)
	if err != nil {
		return err
	}

	// write msg type
	if err := binary.Write(c, binary.LittleEndian, t); err != nil {
		return err
	}

	// write msg length
	if err := binary.Write(c, binary.LittleEndian, uint16(len(b))); err != nil {
		return err
	}
	// write msg
	_, err = c.Write(b)
	return err
}

type Tunnel struct {
	Hash           string
	Type           string
	Destination    string
	DestinationTLS bool

	// use
	ctrlCon    *Con
	child      []*Con
	quit       chan bool
	locker     *sync.Mutex
	PublicAddr string
}

const cloudTunAddr = "tun.vkng.cc:443"

type NetEvent struct {
	action uint8
	data   []byte
}

func (t *Tunnel) Start(errChan chan error) error {
	t.locker = &sync.Mutex{}
	ctrlCon, err := connect(cloudTunAddr, true)
	if err != nil {
		return err
	}

	t.quit = make(chan bool)
	t.ctrlCon = &Con{ctrlCon, &sync.RWMutex{}}

	go func() {
		defer func() {
			ctrlCon.Close()
		}()

		err := t.ctrlCon.WritePacket(AuthReq, &AuthReqMsg{Hash: t.Hash, Type: t.Type, Destination: t.Destination})
		if err != nil {
			utils.PrintDebug("failed to write packet %v", err)
			errChan <- err
			return
		}

		netEventCh := make(chan NetEvent)
		netErrCh := make(chan error)
		go func(netEventCh chan NetEvent, netErrCh chan error) {
			for {
				action, data, err := t.ctrlCon.ReadPacket()
				if err != nil {
					netErrCh <- err
					return
				}
				netEventCh <- NetEvent{action, data}
			}
		}(netEventCh, netErrCh)

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case e := <-netErrCh:
				utils.PrintDebug("netErrCh %v", e)
				close(netErrCh)
				errChan <- e
				return
			case <-t.quit:
				t.ctrlCon.Close()
				t.clear()
				close(t.quit)
				close(errChan)
				return
			case <-ticker.C:
				t.ping()
			case event := <-netEventCh:
				switch event.action {
				case Pong:
				case Stat:
					statMsg := &shared.StatMsg{}
					err = msgpack.Unmarshal(event.data, statMsg)
					if err != nil {
						log.Printf("failed to unmarshal stat %v", err)
						continue
					}
					shared.Stat[t.Hash] = *statMsg
					log.Printf("tunnel %s stat: %v", t.Hash, shared.Stat[t.Hash])
				case AuthRes:
					utils.PrintDebug("got AuthRes %s", string(event.data))
					authResMsg := &AuthResMsg{}
					err = msgpack.Unmarshal(event.data, authResMsg)
					if err != nil {
						utils.PrintDebug("unmarshal err %v", err)
						errChan <- err
						break
					}
					t.PublicAddr = authResMsg.PublicAddr
				case NewConReq:
					utils.PrintDebug("request for new connection")
					newConReqMsg := &NewConReqMsg{}
					err = msgpack.Unmarshal(event.data, newConReqMsg)
					if err != nil {
						errChan <- err
						break
					}
					serverCon, err := connect(cloudTunAddr, true)
					if err != nil {
						errChan <- err
						break
					}

					dstCon, err := connect(t.Destination, t.DestinationTLS)
					if err != nil {
						log.Println("failed to connect to dst", err)
						errChan <- err
						break
					}

					err = serverCon.WritePacket(NewConRes, &NewConResMsg{Session: newConReqMsg.Session})
					if err != nil {
						utils.PrintDebug("Failed to write packet %s", err)
						errChan <- err
						break
					}
					t.addChild(serverCon)
					t.addChild(dstCon)
					go t.link(dstCon, serverCon)
				}
			}
		}
	}()
	return nil
}

func (t *Tunnel) Stop() {
	t.quit <- true
}

func (t *Tunnel) ping() {
	t.ctrlCon.WritePacket(Ping, struct{}{})
}

func (t *Tunnel) link(src, dst *Con) {
	go io.Copy(src, dst)
	io.Copy(dst, src)
	src.Close()
	dst.Close()
}

func (t *Tunnel) addChild(c *Con) {
	t.locker.Lock()
	defer t.locker.Unlock()
	t.child = append(t.child, c)
}

func (t *Tunnel) clear() {
	t.locker.Lock()
	defer t.locker.Unlock()
	for _, c := range t.child {
		if c != nil {
			c.Close()
		}
	}
	t.child = []*Con{}
}

func connect(addr string, useTLS bool) (*Con, error) {
	utils.PrintDebug("new connection to %s; tls:%v", addr, useTLS)
	var (
		con net.Conn
		err error
	)

	if useTLS {
		host, _, _ := net.SplitHostPort(addr)
		con, err = tls.Dial("tcp", addr, &tls.Config{
			ServerName:         host, // set SNI
			InsecureSkipVerify: true, // ignore invalid authority
		})
	} else {
		con, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return nil, err
	}
	return &Con{con, &sync.RWMutex{}}, err
}
