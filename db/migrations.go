package db

import (
	"log"

	"github.com/vikingo-project/vsat/models"
)

func migrate() {
	err := connection.AutoMigrate(&models.Service{}, &models.Session{},
		&models.FullEvent{}, &models.File{}, &models.Crt{}, &models.Auth{})
	if err != nil {
		log.Fatal(err)
	}
	/*
		if utils.IsDevMode() {
			var count int64
			connection.Begin().Model(&models.Service{}).Count(&count)
			if count < 1 {

				var m modules.Module
				m = vstcp.Load()
				mSettings, _ := json.Marshal(m.GetDefaultSettings())
				err = connection.Begin().Save(&models.Service{
					Hash:       utils.NewUUID(),
					ModuleName: "TCP",
					ListenIP:   "0.0.0.0",
					ListenPort: 8081,
					Autostart:  false,
					Settings:   string(mSettings),
				})
				if err != nil {
					log.Println(err)
				}

				m = vshttp.Load()
				mSettings, _ = json.Marshal(m.GetDefaultSettings())
				err = connection.Begin().Save(&models.Service{
					Hash:       utils.NewUUID(),
					ModuleName: "HTTP",
					ListenIP:   "0.0.0.0",
					ListenPort: 8082,
					Autostart:  false,
					Settings:   string(mSettings),
				})
				if err != nil {
					log.Println(err)
				}
			}
			connection.Begin().Model(&models.File{}).Count(&count)
			if count < 1 {
				err = connection.Begin().Save(&models.File{
					Hash:        utils.NewUUID(),
					FileName:    "Example",
					ContentType: "text/plain",
					Size:        6,
					Data:        []byte("kekmek"),
				})

				if err != nil {
					log.Println(err)
				}
			}

			connection.Begin().Model(&models.Event{}).Count(&count)
			if count < 1 {
				err = connection.Begin().Save(&models.Event{
					Hash:       utils.NewUUID(),
					ModuleName: "DNS",
					Service:    "dev-service-hash",
					ClientIP:   "127.0.0.1",
					Data:       `{"resource":"vikingo.org","type":"A"}`,
					Extra:      `{"response":"127.0.0.1"}`,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	*/
}
