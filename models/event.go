package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	Hash        string    `json:"hash" gorm:"primaryKey"`
	Date        time.Time `json:"date"`
	ClientIP    string    `json:"client_ip" gorm:"index"`
	Description string    `json:"description" gorm:"index"`
	ModuleName  string    `json:"module_name"`
	Service     string    `json:"service" gorm:"index"` // hash
	ServiceName string    `json:"service_name"`
	Visited     bool      `json:"visited"`
	LocalAddr   string    `json:"local_addr"`
}

type SessionInfo struct {
	Description string `json:"description"`
	LocalAddr   string `json:"local_addr"`
	ClientIP    string `json:"client_ip"`
}

type FullEvent struct {
	gorm.Model
	Hash    string    `json:"hash" gorm:"primaryKey"`
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
