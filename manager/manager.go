package manager

import (
	"log"
	"time"

	"github.com/vikingo-project/vsat/modules"
	vsdns "github.com/vikingo-project/vsat/modules/dns"
	vsftp "github.com/vikingo-project/vsat/modules/ftp"
	vshttp "github.com/vikingo-project/vsat/modules/http"
	vsroguemysql "github.com/vikingo-project/vsat/modules/rogue-mysql"
	vstcp "github.com/vikingo-project/vsat/modules/tcp"

	"github.com/vikingo-project/vsat/utils"
)

type InstanceInfo struct {
	Module  *modules.Module
	Started time.Time
}

type Manager struct {
	Instances map[string]InstanceInfo
}

func (mgr *Manager) Start() {
	var module modules.Module

	module = vsdns.Load()
	modules.Register(module)

	module = vshttp.Load()
	modules.Register(module)

	module = vstcp.Load()
	modules.Register(module)

	module = vsftp.Load()
	modules.Register(module)

	module = vsroguemysql.Load()
	modules.Register(module)

	// check autostarted services in the DB
	services, _ := loadServicesFromDB()
	for _, service := range services {
		if service.Autostart {
			m, err := mgr.startService(service)
			if err != nil {
				log.Printf("Failed to start service %v", err)
				continue
			}

			mgr.Instances[service.Hash] = InstanceInfo{
				Module:  &m,
				Started: time.Now(),
			}
		}
	}

	// start services with Autostart true
}

func (mgr *Manager) StartService(hash string) error {
	// todo: check running instances by hash...
	service, err := loadServiceFromDB(hash)
	if err != nil {
		return err
	}
	if utils.IsDevMode() {
		log.Printf("manager: service %s loaded from db", service.Hash)
	}

	m, err := mgr.startService(service)
	if err != nil {
		log.Printf("Failed to start service (%s)", err)
		return err
	}

	mgr.Instances[service.Hash] = InstanceInfo{
		Module:  &m,
		Started: time.Now(),
	}

	return nil
}

func (mgr *Manager) StopService(hash string) error {
	return mgr.stopService(hash)
}

func (mgr *Manager) IsServiceActive(hash string) bool {
	if _, ok := mgr.Instances[hash]; ok {
		return true
	}
	return false
}

func NewManager() *Manager {
	return &Manager{
		Instances: make(map[string]InstanceInfo),
	}
}
