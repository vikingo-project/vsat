package models

import "github.com/akamajoris/ngorm/model"

type Auth struct {
	model.Model
	Token string `json:"token"`
}
