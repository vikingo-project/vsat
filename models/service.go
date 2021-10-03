package models

import "github.com/akamajoris/ngorm/model"

type Service struct {
	model.Model
	Hash        string `json:"hash"`
	ServiceName string `json:"serviceName"`
	ModuleName  string `json:"moduleName"`
	ListenIP    string `json:"listenIP"`
	ListenPort  int    `json:"listenPort"`
	Autostart   bool   `json:"autostart"`
	Active      bool   `json:"active"`
	Settings    string `json:"settings"`
	BaseProto   string `json:"baseProto"`
}
