package models

import (
	"time"

	"github.com/akamajoris/ngorm/model"
)

type File struct {
	model.Model
	Date            time.Time `mapstructure:"date" json:"date"`
	Hash            string    `mapstructure:"hash" json:"hash"`
	FileName        string    `mapstructure:"file_name" json:"file_name"`
	ContentType     string    `mapstructure:"content_type" json:"content_type"`
	Size            int64     `mapstructure:"size" json:"size"`
	Data            []byte    `mapstructure:"data" json:"data"`
	InteractionHash string    `mapstructure:"interaction_hash" json:"interaction_hash"`
	ServiceHash     string    `mapstructure:"service_hash" json:"service_hash"`
}
