package events

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/shared"
	"github.com/vikingo-project/vsat/utils"
)

func NewAPI(serviceHash, serviceName, moduleName string) *EventsAPI {
	return &EventsAPI{
		serviceHash: serviceHash,
		serviceName: serviceName,
		moduleName:  moduleName,
	}
}

type EventsAPI struct {
	moduleName  string
	serviceHash string
	serviceName string
}

func (api *EventsAPI) NewSession(info models.SessionInfo) (string, error) {
	hash := utils.EasyHash(true)
	err := db.GetConnection().Begin().Save(&models.Session{
		Hash:        hash,
		Date:        time.Now(),
		Service:     api.serviceHash,
		ServiceName: api.serviceName,
		ModuleName:  api.moduleName,
		ClientIP:    info.ClientIP,
		Description: info.Description, // short info
		LocalAddr:   info.LocalAddr,
	})
	if err != nil {
		log.Println("failed to start new session", err)
		return "", err
	}
	shared.Updates <- struct {
		Name string `json:"name"`
	}{"interactions"}

	return hash, nil
}

func (api *EventsAPI) PushEvent(event models.Event) {
	// store files
	for k, v := range event.Data {
		if f, ok := v.(models.File); ok {
			fileHash := utils.EasyHash(true)
			db.GetConnection().Begin().Save(&f)
			delete(event.Data, k)
			event.Data["file:"+k] = fileHash
		}
	}

	encodedData, _ := json.Marshal(event.Data)
	err := db.GetConnection().Begin().Save(&models.FullEvent{
		Hash:    utils.EasyHash(true),
		Date:    time.Now(),
		Name:    event.Name,
		Session: event.Session,
		Data:    string(encodedData),
	})
	if err != nil {
		log.Println(err)
	}

	// update file fileds
	for k, v := range event.Data {
		if strings.HasPrefix(k, "file:") {
			fileHash := v.(string)
			db.GetConnection().Model(&models.File{}).Where("hash = ?", fileHash).Update("interaction_hash", event.Session)
		}
	}

	shared.Updates <- struct {
		Name string `json:"name"`
	}{"event"}
}
