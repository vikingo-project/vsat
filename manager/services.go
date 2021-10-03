package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/events"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/modules"
	"github.com/vikingo-project/vsat/utils"
)

func loadServicesFromDB() ([]models.Service, error) {
	var services []models.Service
	err := db.GetConnection().Begin().Find(&services)
	return services, err
}

func loadServiceFromDB(hash string) (models.Service, error) {
	var (
		service = models.Service{Hash: hash}
		count   int64
	)

	db.GetConnection().Begin().Model(&service).Where(&service).Count(&count)
	if count != 1 {
		return service, errors.New("service not found in DB")
	}

	err := db.GetConnection().Begin().Find(&service, &service)
	return service, err
}

func (mgr *Manager) startService(service models.Service) (modules.Module, error) {
	// get data from db by hash
	module := modules.GetModuleByName(service.ModuleName)
	if module == nil {
		return nil, fmt.Errorf("module %s is not registered", service.ModuleName)
	}
	if utils.IsDevMode() {
		log.Println("Start service", service.Settings)
	}

	// decode settings from DB (string) to map[string]interface{}
	var settings interface{}
	err := json.Unmarshal([]byte(service.Settings), &settings) // error doesn't matter ;)
	if err != nil {
		log.Println("Failed to decode JSON from DB", err)
	}

	API := events.NewAPI(service.Hash, service.ServiceName, service.ModuleName)
	module.Init(service.ListenIP, service.ListenPort, settings, API)

	var (
		ticker  = time.NewTicker(time.Second)
		errChan = make(chan error)
	)
	go func(errChan chan error) {
		errChan <- module.Up()
	}(errChan)
	select {
	// waiting for error or timeout
	case err := <-errChan:
		log.Println("Got error on start", err)
		return module, err
	case <-ticker.C:
		log.Println("err timeout")
		return module, nil
	}

}

func (mgr *Manager) stopService(hash string) error {
	if m, ok := mgr.Instances[hash]; ok {
		module := *m.Module
		err := module.Down()
		if err != nil {
			log.Println("Failed to stop service", err.Error())
			return err
		}
		delete(mgr.Instances, hash)
	}
	log.Println("service has been stopped")
	return nil
}
