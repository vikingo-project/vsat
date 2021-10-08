package models

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	Token string `json:"token"`
}
