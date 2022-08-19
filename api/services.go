package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/models"
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
