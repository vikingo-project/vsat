package shared

type config struct {
	DB     string `yaml:"db"`               // use persistent storage (sqlite3); default: false
	Listen string `yaml:"ctrl_listen_addr"` // default: '0.0.0.0:3419'
	Debug  bool   `yaml:"debug"`            // enable debug mode; default: false
	Token  string `yaml:"token"`            // ctrl access token
}

var (
	Config    config
	Version   = "0.0.0" // the real versions gets from build env
	BuildHash = "-"
	Updates   = make(chan interface{})
)
