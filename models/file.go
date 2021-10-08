package models

import (
	"time"
)

type File struct {
	Hash            string    `json:"hash" gorm:"primaryKey"`
	Date            time.Time `json:"date"`
	FileName        string    `json:"file_name"`
	ContentType     string    `json:"content_type"`
	Size            int64     `json:"size"`
	Data            []byte    `json:"data"`
	InteractionHash string    `json:"interaction_hash"`
	ServiceHash     string    `json:"service_hash"`
}
