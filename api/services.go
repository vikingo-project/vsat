package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/modules"
	"github.com/vikingo-project/vsat/utils"
	"gorm.io/gorm"
)

func (a *APIC) Services(params string) (*RecordsContainer, error) {
	q, _ := url.ParseQuery(params)
	filterBaseProto := strings.Trim(q.Get("base_proto"), " ")
	filterName := strings.TrimSpace(q.Get("service_name"))

	var services []models.Service
	dq := db.GetConnection().Model(&models.Service{})
	if filterBaseProto != "" {
		dq.Where(`base_proto LIKE ?`, fmt.Sprintf("%%%s%%", strings.ToLower(filterBaseProto)))
	}

	if filterName != "" {
		dq.Where("service_name LIKE ?", fmt.Sprintf("%%%s%%", filterName))
	}

	err := dq.Find(&services).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &RecordsContainer{}, err
	}

	for i, service := range services {
		services[i].Active = manager.IsServiceActive(service.Hash) // set actual service status
	}
	return &RecordsContainer{
		Records: services,
	}, nil
}

func (a *APIC) CreateService(ws *models.WebService) (string, error) {
	vdtr := validator.New()
	vdtr.SetTagName("binding")

	// todo: validate settings...
	if err := vdtr.Struct(ws); err != nil {
		return "", err
	}
	module := modules.GetModuleByName(ws.ModuleName)
	service := &models.Service{
		Hash:        utils.EasyHash(false),
		ServiceName: ws.ServiceName,
		ModuleName:  ws.ModuleName,
		ListenIP:    ws.ListenIP,
		ListenPort:  ws.ListenPort,
		Autostart:   ws.Autostart,
		BaseProto:   strings.Join(module.GetInfo()["base_proto"].([]string), "/"),
		Settings:    []byte(ws.Settings),
	}

	err := db.GetConnection().Save(service).Error
	return service.Hash, err
}

func (a *APIC) UpdateService(params *models.WebService) (string, error) {
	vdtr := validator.New()
	vdtr.SetTagName("binding")

	// todo: validate settings...
	if err := vdtr.Struct(params); err != nil {
		return "", err
	}
	err := db.GetConnection().Model(&models.Service{}).Where(&models.Service{Hash: params.Hash}).Updates(&models.Service{ServiceName: params.ServiceName,
		ListenIP: params.ListenIP, ListenPort: params.ListenPort, Autostart: params.Autostart, Settings: []byte(params.Settings)}).Error

	return params.Hash, err
}

func (a *APIC) ToggleService(params *models.ChangeServiceState) error {
	vdtr := validator.New()
	vdtr.SetTagName("binding")
	if err := vdtr.Struct(params); err != nil {
		return err
	}
	if params.NewState == "start" {
		return manager.StartService(params.Hash)

	} else {
		return manager.StopService(params.Hash)
	}
}
