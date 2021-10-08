package models

import "gorm.io/gorm"

// single certificate
type Crt struct {
	gorm.Model
	Name string `mapstructure:"name" json:"name"`
	Data []byte `mapstructure:"data" json:"data"`
}
