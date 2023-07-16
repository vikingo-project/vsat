package vshttp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
	"golang.org/x/crypto/acme/autocert"
)

type header struct {
	Name  string `json:"name" mapstructure:"name"`
	Value string `json:"value" mapstructure:"value"`
}
type templateSettings struct {
	Status   int      `json:"status" mapstructure:"status"`
	Template string   `json:"template" mapstructure:"template"`
	Headers  []header `json:"headers" mapstructure:"headers"`
}
type fileSettings struct {
	Hash string `json:"hash" mapstructure:"hash"`
}

type folderSettings struct {
	Folder string `json:"folder" mapstructure:"folder"`
}

type proxySettings struct {
	Destination string `json:"destination" mapstructure:"destination"`
	CustomHost  string `json:"custom_host" mapstructure:"custom_host"`
}

type Location struct {
	Path   string      `json:"path" mapstructure:"path"`               // location path/URI
	Action string      `json:"action_name" mapstructure:"action_name"` // what to do
	Data   interface{} `json:"action_data" mapstructure:"action_data"` // extra data: template/file id/location/etc
}

type HostConfig struct {
	Hostname  string     `json:"hostname" mapstructure:"hostname"`   // slice of hostnames
	Locations []Location `json:"locations" mapstructure:"locations"` // slice of locations with settings
}

type TLSConfig struct {
	Enabled  bool `json:"enabled" mapstructure:"enabled"`
	Autocert bool `json:"autocert" mapstructure:"autocert"`
}

type settings struct {
	TLS             TLSConfig    `json:"tls" mapstructure:"tls"`
	Hosts           []HostConfig `json:"hosts" mapstructure:"hosts"`
	AllowFileUpload bool         `json:"allow_file_upload" mapstructure:"allow_file_upload"`
	LogRequest      bool         `json:"log_request" mapstructure:"log_request"`
	LogResponse     bool         `json:"log_response" mapstructure:"log_response"`
}

type module struct {
	Name        string
	Description string
	BaseProto   []string

	listenIP   string
	listenPort int
	settings   settings
	server     *http.Server
	API        *events.EventsAPI
}

func Load() *module {
	return &module{
		Name:        "HTTP",
		Description: "HTTP/S server.",
		BaseProto:   []string{"tcp"},
	}
}

func (m *module) Init(listenIP string, listenPort int, settings interface{}, API *events.EventsAPI) {
	m.listenIP = listenIP
	m.listenPort = listenPort
	m.API = API
	log.Println("Init HTTP with", m.settings, settings)
	err := utils.ExtractSettings(&m.settings, settings)
	if err != nil {
		log.Println("Failed to decode module settings to struct", err)
	}
}

func (m *module) GetName() string {
	return m.Name
}

func (m *module) GetDefaultSettings() interface{} {
	s := &settings{
		Hosts: []HostConfig{
			{
				Hostname: "*",
				Locations: []Location{
					{Path: "/", Action: "template", Data: "hi there"},
					{Path: "/test.exe", Action: "file", Data: "@test"},
				},
			},
		},
	}
	return s
}

func (m *module) GetSettings() interface{} {
	return m.settings
}

func (m *module) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"base_proto":  m.BaseProto,
	}
}

func (m *module) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fakeWriter := httptest.NewRecorder()
		next.ServeHTTP(fakeWriter, r)

		// send response to the client
		for name, values := range fakeWriter.Header() {
			for _, v := range values {
				w.Header().Add(name, v)
			}
		}
		w.WriteHeader(fakeWriter.Code) // shold be sent after headers!
		w.Write(fakeWriter.Body.Bytes())

		clientIP := utils.ExtractIP(r.RemoteAddr)

		session, _ := m.API.NewSession(models.SessionInfo{
			Description: fmt.Sprintf("%s %s", r.Method, r.RequestURI),
			ClientIP:    clientIP,
			LocalAddr:   fmt.Sprintf("%s:%d", m.listenIP, m.listenPort),
		})

		var data = make(map[string]interface{})

		if m.settings.LogRequest {
			rawReq, _ := httputil.DumpRequest(r, true)
			// if len(rawReq) > 1024*10 {
			// 	rawReq = rawReq[:1024*10]
			// 	rawReq = append(rawReq, []byte("...")...)
			// }
			data["request"] = string(rawReq)
		}

		if m.settings.LogResponse {
			rawRes, _ := httputil.DumpResponse(fakeWriter.Result(), true)
			// if len(rawRes) > 1024*10 {
			// 	rawRes = rawRes[:1024*10]
			// 	rawRes = append(rawRes, []byte("...")...)
			// }
			data["response"] = string(rawRes)
		}

		if m.settings.LogRequest || m.settings.LogResponse {
			m.API.PushEvent(models.Event{Session: session, Name: "request", Data: data})
		}

		if m.settings.AllowFileUpload {
			err := r.ParseMultipartForm(10 << 20)
			if err == nil {
				formdata := r.MultipartForm
				for k, v := range formdata.File {
					log.Println("got kv", k, v)
					for _, f := range v {
						name := f.Filename
						size := f.Size
						// 50m
						if size > 1024*1024*50 {
							size = 1024 * 1024 * 50
						}

						file, ferr := f.Open()
						if ferr == nil {
							buff := make([]byte, size)
							_, rerr := file.Read(buff)
							if rerr == nil {
								fileStruct := files.PrepareFile(name, buff)
								m.API.PushEvent(models.Event{Session: session, Name: "File upload", Data: map[string]interface{}{
									"form key": k,
									"filename": fileStruct.FileName,
									"file":     fileStruct,
								}})
							} else {
								// ignore error
							}
							file.Close()
						}

					}
				}
			}
		}

	})
}

type hostnames struct {
	mu   sync.Mutex
	list map[string]bool
}

func (h *hostnames) Add(host string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.list[host] = true
}

func (h *hostnames) GetHostnames() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	list := []string{}
	for hostname := range h.list {
		list = append(list, hostname)
	}
	return list
}
func serveTemplate(path, action string, data interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("custom handler", r.URL.Path, r.Header.Get("Host"))
		var settings templateSettings
		err := mapstructure.Decode(data, &settings)
		if err != nil {
			log.Println("Failed to decode template settings")
		}
		tpl := ""
		data, err := Render(settings.Template, r)
		if err != nil {
			log.Println(err)
			data = []byte(tpl) // if there is an error return raw template
		}
		for _, header := range settings.Headers {
			w.Header().Add(header.Name, header.Value)
		}
		w.Write([]byte(data))
	}
}
func serveFile(path, action string, data interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("serve file")
		var settings fileSettings
		err := mapstructure.Decode(data, &settings)
		if err != nil {
			log.Println("Failed to decode template settings")
		}

		var file models.File
		err = db.GetConnection().Select("file_name,data,content_type").Where(&models.File{Hash: settings.Hash}).Find(&file).Error
		if err != nil {
			log.Println("fuck", err)
		}

		w.Header().Add("Content-Type", file.ContentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(file.Data)))
		w.Write([]byte(file.Data))
	}
}

func serveFolder(path, action string, data interface{}) http.Handler {
	var settings folderSettings
	err := mapstructure.Decode(data, &settings)
	if err != nil {
		log.Println("Failed to decode template settings")
	}
	folder := strings.TrimRight(settings.Folder, "/")
	return http.FileServer(http.Dir(folder))
}

func doProxy(path, action string, data interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings proxySettings
		err := mapstructure.Decode(data, &settings)
		if err != nil {
			log.Println("Failed to decode template settings")
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr, CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		newURL, _ := url.Parse(settings.Destination)
		req, err := http.NewRequest(r.Method, newURL.Scheme+"://"+newURL.Host+r.RequestURI, r.Body)
		if err == nil {
			req.Header = r.Header
			req.ContentLength = r.ContentLength
			req.Form = r.Form
			req.Method = r.Method
			req.MultipartForm = r.MultipartForm
			req.PostForm = r.PostForm
			req.Proto = r.Proto
			req.ProtoMajor = r.ProtoMajor
			req.ProtoMinor = r.ProtoMinor
			req.Body = r.Body
			req.Close = r.Close

			// req.TLS = r.TLS
			// req.Trailer = r.Trailer
			// req.TransferEncoding = r.TransferEncoding

			if settings.CustomHost != "" {
				r.Host = settings.CustomHost // update old r because http recorder uses it
				req.Host = settings.CustomHost
			} else {
				req.Host = r.Host
			}

			resp, err := client.Do(req)
			if err == nil {
				for name, values := range resp.Header {
					for _, v := range values {
						w.Header().Add(name, v)
					}
				}
				w.WriteHeader(resp.StatusCode)
				body, _ := ioutil.ReadAll(resp.Body)
				w.Write(body)
				defer resp.Body.Close()

			} else {
				// response error
				log.Println("failed to get reponse", err)
			}
		} else {
			// new request error
		}

	}
}

func (m *module) Up() error {
	utils.PrintDebug("Starting HTTP server...")
	router := mux.NewRouter()
	router.NotFoundHandler = m.loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hi")) }))
	router.Use(m.loggingMiddleware)

	thereIsCharHostname := false
	hosts := hostnames{
		mu:   sync.Mutex{},
		list: make(map[string]bool),
	}
	// register matchers
	for _, host := range m.settings.Hosts {
		hosts.Add(host.Hostname)
		for _, location := range host.Locations {

			// check if one of hostname is hostname; if hostname == * -> will be added to the router last
			if host.Hostname == "*" {
				thereIsCharHostname = true
			} else {
				// check if hostname already contains port; for ex.: localhost:8080
				if (m.listenPort != 80 && m.listenPort != 443) && len(strings.Split(host.Hostname, ":")) < 2 {
					host.Hostname = fmt.Sprintf("%s:%d", host.Hostname, m.listenPort)
				}

				func(path, action string, data interface{}) {
					if action == "template" {
						router.Host(host.Hostname).Path(path).HandlerFunc(serveTemplate(path, action, data))
					}
					if action == "file" {
						router.Host(host.Hostname).Path(path).HandlerFunc(serveFile(path, action, data))
					}
					if action == "folder" {
						router.PathPrefix(path).Handler(http.StripPrefix(path, serveFolder(path, action, data)))
					}
					if action == "proxy" {
						router.Host(host.Hostname).PathPrefix(path).HandlerFunc(doProxy(path, action, data))
					}

				}(location.Path, location.Action, location.Data)
			}
		}

	}

	// * hostname should be added last; todo: rewrite mux Matcher
	if thereIsCharHostname {
		for _, host := range m.settings.Hosts {
			if host.Hostname == "*" {
				for _, location := range host.Locations {
					func(path, action string, data interface{}) {
						if action == "template" {
							router.Path(path).HandlerFunc(serveTemplate(path, action, data))
						}
						if action == "file" {
							router.Path(path).HandlerFunc(serveFile(path, action, data))
						}
						if action == "folder" {
							router.PathPrefix(path).Handler(http.StripPrefix(path, serveFolder(path, action, data)))
							// router.Path(path).Handler(serveFolder(path, action, data))
						}
						if action == "proxy" {
							router.PathPrefix(path).HandlerFunc(doProxy(path, action, data))
						}
					}(location.Path, location.Action, location.Data)
				}
			}
		}
	}

	srv := &http.Server{}
	useTLS := m.settings.TLS.Enabled
	if useTLS {
		var tlsConf *tls.Config
		if m.settings.TLS.Autocert {
			cc := CertsCache{}
			am := &autocert.Manager{
				Cache: cc,
				// Cache:      autocert.DirCache("/tmp"),
				Prompt: autocert.AcceptTOS,
				HostPolicy: func(ctx context.Context, host string) error {
					hostsList := hosts.GetHostnames()
					// sort by length desc
					sort.Slice(hostsList, func(i, j int) bool {
						return len(hostsList[i]) > len(hostsList[j])
					})
					// check subdomains, wildcards
					for _, h := range hostsList {
						if match(h, host) {
							return nil
						}
					}
					return fmt.Errorf("acme/autocert: host %q not configured in HostWhitelist", host)
				},
			}
			tlsConf = am.TLSConfig()
		} else {
			// todo: get it from cert store
			crt, err := utils.GenerateCert("test.com")
			if err != nil {
				log.Println("Failed to load cert", err)
			}
			tlsConf = &tls.Config{Certificates: []tls.Certificate{crt}}
		}

		srv = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", m.listenIP, m.listenPort),
			Handler:      router,
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
			TLSConfig:    tlsConf,
		}
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Printf("listen: %s\n", err)
			}
		}()

	} else {
		srv = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", m.listenIP, m.listenPort),
			Handler:      router,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Printf("listen: %s\n", err)
			}
		}()
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below

	m.server = srv

	return nil
}

func (m *module) Down() error {
	log.Println("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	return nil
}

// match returns true if wildcard value is mathed
func match(name string, value string) bool {
	var result strings.Builder
	for i, literal := range strings.Split(name, "*") {
		if i > 0 {
			result.WriteString(".*")
		}
		result.WriteString(regexp.QuoteMeta(literal))
	}
	matched, _ := regexp.MatchString(result.String(), value)
	return matched
}
