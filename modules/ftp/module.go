package vsftp

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/utils"
)

type fsn struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Dir   bool   `json:"dir"`
	Size  int64  `json:"size"`
}

type settings struct {
	MinPassivePort int            `json:"minassivePort"`
	MaxPassivePort int            `json:"maxPassivePort"`
	FS             map[string]fsn `json:"fs"`
}

type module struct {
	Name        string
	Description string
	BaseProto   []string

	listenIP   string
	listenPort int
	settings   settings
	API        *events.EventsAPI
	Server     *Server
}

func Load() *module {
	return &module{
		Name:        "FTP",
		Description: "FTP Server",
		BaseProto:   []string{"tcp"},
	}
}

func copyVFS(m sync.Map) sync.Map {
	var n sync.Map
	m.Range(func(k, v interface{}) bool {
		vm, ok := v.(sync.Map)
		if ok {
			n.Store(k, copyVFS(vm))
		} else {
			n.Store(k, v)
		}
		return true
	})
	return n
}

func (m *module) Init(listenIP string, listenPort int, settings interface{}, API *events.EventsAPI) {
	m.listenIP = listenIP
	m.listenPort = listenPort
	m.API = API
	err := utils.ExtractSettings(&m.settings, settings)
	log.Println("decode err", err)

	fsNodes := sync.Map{}
	for path, node := range m.settings.FS {
		size := int64(0)
		if !node.Dir {
			// get file info
			file := files.GetFileByHash(node.Value, false)
			size = file.Size
		}
		fsNodes.Store(path, &vfsNode{name: node.Name, dir: node.Dir, size: size, value: node.Value})
	}
	/*
		fsNodes.Store("/", &vfsNode{name: ".", dir: true, size: 0})
		fsNodes.Store("/test", &vfsNode{name: "test", dir: false, size: 150})
	*/

	m.Server = &Server{
		EventsAPI: API,
		settings:  m.settings,
		vfs:       VirtualFS{nodes: fsNodes},
		quit:      make(chan interface{}),
	}
}

func (m *module) GetName() string {
	return m.Name
}

func (m *module) GetSettings() interface{} {
	return m.settings
}

func (m *module) GetDefaultSettings() interface{} {
	return settings{}
}

func (m *module) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"base_proto":  m.BaseProto,
	}
}

// listenIP, listenPort, settings, pusher
func (m *module) Up() error {
	l, err := net.Listen("tcp", net.JoinHostPort(m.listenIP, strconv.Itoa(m.listenPort)))
	if err != nil {
		return err
	}
	m.Server.listener = l
	m.Server.wg.Add(1)
	go m.Server.serve()
	return nil
}

func (m *module) Down() error {
	m.Server.Stop()
	return nil
}
