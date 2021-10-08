package models

type Service struct {
	Hash        string `json:"hash" gorm:"primary_key"`
	ServiceName string `json:"serviceName"`
	ModuleName  string `json:"moduleName"`
	ListenIP    string `json:"listenIP"`
	ListenPort  int    `json:"listenPort"`
	Autostart   bool   `json:"autoStart"`
	Active      bool   `json:"active"`
	Settings    string `json:"settings"`
	BaseProto   string `json:"baseProto"`
}
