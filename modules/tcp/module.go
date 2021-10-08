package vstcp

import (
	"net"
	"strconv"

	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/utils"
)

type ResponseSettings struct {
	Response string `json:"response" mapstructure:"response"`
	// todo: add steps/matchers
}

type ProxySettings struct {
	Destination string `json:"destination" mapstructure:"destination"`
}

type settings struct {
	LogRequest       bool             `json:"log_request" mapstructure:"log_request"`
	LogResponse      bool             `json:"log_response" mapstructure:"log_response"`
	Mode             string           `json:"mode" mapstructure:"mode"`
	ProxySettings    ProxySettings    `json:"proxy_settings" mapstructure:"proxy_settings"`
	ResponseSettings ResponseSettings `json:"response_settings" mapstructure:"response_settings"`
}

type module struct {
	Name        string
	Description string
	BaseProto   []string

	listenIP   string
	listenPort int
	settings   settings

	API    events.EventsAPI
	Server *Server
}

func Load() *module {
	return &module{
		Name:        "TCP",
		Description: "Simple TCP server with proxy mode",
		BaseProto:   []string{"tcp"},
	}
}

func (m *module) Init(listenIP string, listenPort int, settings interface{}, API *events.EventsAPI) {
	m.listenIP = listenIP
	m.listenPort = listenPort
	m.API = *API
	utils.ExtractSettings(&m.settings, settings)

	m.Server = &Server{
		quit:     make(chan interface{}),
		API:      API,
		settings: m.settings,
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
	// todo: "tcp" is listening ipv6 interfaces only in docker
	l, err := net.Listen("tcp4", net.JoinHostPort(m.listenIP, strconv.Itoa(m.listenPort)))
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
