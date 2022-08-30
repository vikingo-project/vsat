package models

import "gorm.io/datatypes"

type Service struct {
	Hash        string         `json:"hash" gorm:"primary_key"`
	ServiceName string         `json:"serviceName"`
	ModuleName  string         `json:"moduleName"`
	ListenIP    string         `json:"listenIP"`
	ListenPort  int            `json:"listenPort"`
	Settings    datatypes.JSON `json:"moduleSettings"`
	Autostart   bool           `json:"autoStart"`
	Active      bool           `json:"active"`
	BaseProto   string         `json:"baseProto"`
}

type WebService struct {
	Hash        string      `json:"hash" `
	ServiceName string      `json:"serviceName" binding:"required"`
	ModuleName  string      `json:"moduleName" binding:"required"`
	ListenIP    string      `json:"listenIP" binding:"required"`
	ListenPort  int         `json:"listenPort" binding:"required"`
	Settings    interface{} `json:"moduleSettings" binding:"required"`
	Autostart   bool        `json:"autoStart"`
	Active      bool        `json:"active"`
	BaseProto   string      `json:"baseProto"`
}

type ServiceHash struct {
	Hash string `json:"hash" binding:"required"`
}

type ChangeServiceState struct {
	Hash     string `json:"hash" binding:"required"`
	NewState string `json:"state" binding:"required"`
}
