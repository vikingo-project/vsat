package vsdns

import (
	"time"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/utils"
)

type settings struct {
	Recursive bool     `json:"recursive"` // enable recursive requests
	Records   []Record `json:"records"`   // DNS records
	Resolvers []string `json:"resolvers"` // resolv nameservers
}

type module struct {
	Name        string
	Description string
	BaseProto   []string

	listenIP   string
	listenPort int

	Settings settings

	server *Server
	API    *events.EventsAPI
}

func Load() *module {
	return &module{
		Name:        "DNS",
		Description: "Simple DNS server",
		BaseProto:   []string{"tcp", "udp"},
		Settings:    settings{},
	}
}

func (m *module) Init(listenIP string, listenPort int, settings interface{}, API *events.EventsAPI) {
	m.listenIP = listenIP
	m.listenPort = listenPort
	m.API = API
	utils.ExtractSettings(&m.Settings, settings)
}

func (m *module) GetName() string {
	return m.Name
}

func (m *module) GetSettings() interface{} {
	return m.Settings
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

func (m *module) Up() error {
	server := &Server{
		host:     m.listenIP,
		port:     m.listenPort,
		settings: m.Settings,
		API:      m.API, // events API

		rTimeout: 5 * time.Second,
		wTimeout: 5 * time.Second,
	}
	m.server = server
	server.Run()
	return nil
}

func (m *module) Down() error {
	err := m.server.Stop()
	m.server = nil
	return err
}
