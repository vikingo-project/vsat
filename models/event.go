package models

import (
	"time"

	"github.com/akamajoris/ngorm/model"
)

type Session struct {
	model.Model
	Hash        string    `json:"hash"`
	Date        time.Time `json:"date"`
	ClientIP    string    `json:"client_ip"`
	Description string    `json:"description"`
	ModuleName  string    `json:"module_name"`
	Service     string    `json:"service"`      // hash
	ServiceName string    `json:"service_name"` // hash
	Visited     bool      `json:"visited"`
	LocalAddr   string    `json:"local_addr"`
}

type SessionInfo struct {
	Description string `json:"description"`
	LocalAddr   string `json:"local_addr"`
	ClientIP    string `json:"client_ip"`
}

type FullEvent struct {
	model.Model
	Hash    string    `json:"hash"`
	Date    time.Time `json:"date"`
	Session string    `json:"session"`
	Name    string    `json:"name"`
	Data    string    `json:"data"` // JSON encoded fields
}

type Event struct {
	Session string                 `json:"session"` // session_hash
	Name    string                 `json:"name"`
	Data    map[string]interface{} `json:"data"` // JSON encoded fields
}
