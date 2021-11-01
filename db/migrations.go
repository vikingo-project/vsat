package db

import (
	"log"

	"github.com/vikingo-project/vsat/models"
)

func migrate() {
	err := connection.AutoMigrate(&models.Service{}, &models.Session{},
		&models.FullEvent{}, &models.File{}, &models.Crt{}, &models.Auth{}, &models.Tunnel{})
	if err != nil {
		log.Fatal(err)
	}
}
