package models

import (
	"github.com/akamajoris/ngorm/model"
)

// single certificate
type Crt struct {
	model.Model
	Name string `mapstructure:"name" json:"name"`
	Data []byte `mapstructure:"data" json:"data"`
}
