package manager

import (
	"log"
	"sync"
	"time"

	"github.com/vikingo-project/vsat/modules"
	vsdns "github.com/vikingo-project/vsat/modules/dns"
	vsftp "github.com/vikingo-project/vsat/modules/ftp"
	vshttp "github.com/vikingo-project/vsat/modules/http"
	vsroguemysql "github.com/vikingo-project/vsat/modules/rogue-mysql"
	vstcp "github.com/vikingo-project/vsat/modules/tcp"
	"github.com/vikingo-project/vsat/tunnels"

	"github.com/vikingo-project/vsat/utils"
)

type InstanceInfo struct {
	Module  *modules.Module
	Started time.Time
}

type tunnelsList struct {
	locker sync.Mutex
	list   map[string]*tunnels.Tunnel
}

func (tl *tunnelsList) Exists(hash string) bool {
	tl.locker.Lock()
	defer tl.locker.Unlock()
	if _, ok := tl.list[hash]; ok {
		return true
	}
	return false
}

func (tl *tunnelsList) GetPublicAddr(hash string) string {
	tl.locker.Lock()
	defer tl.locker.Unlock()
	if t, ok := tl.list[hash]; ok {
		if t != nil {
			return t.PublicAddr
		}
	}
	return ""
}

func (tl *tunnelsList) Get(hash string) *tunnels.Tunnel {
	tl.locker.Lock()
	defer tl.locker.Unlock()
	return tl.list[hash]
}

func (tl *tunnelsList) Add(hash string, t *tunnels.Tunnel) {
	tl.locker.Lock()
	defer tl.locker.Unlock()
	tl.list[hash] = t
}

func (tl *tunnelsList) Remove(hash string) {
	tl.locker.Lock()
	defer tl.locker.Unlock()
	if t, ok := tl.list[hash]; ok {
		if t != nil {
			delete(tl.list, hash)
		}
	}
}

type Manager struct {
	Instances map[string]InstanceInfo
	Tunnels   tunnelsList
}

func NewManager() *Manager {
	return &Manager{
		Instances: make(map[string]InstanceInfo),
		Tunnels:   tunnelsList{locker: sync.Mutex{}, list: make(map[string]*tunnels.Tunnel)},
	}
}

func (mgr *Manager) Start() {
	// services part
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

	// check services with autostart
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

	// tunnels part
	tunnels, _ := loadTunnelsFromDB()
	utils.PrintDebug("start tunnels %v", tunnels)
	for _, tunnel := range tunnels {
		if tunnel.Autostart {
			err := mgr.startTunnel(tunnel)
			if err != nil {
				log.Printf("Failed to start service %v", err)
				continue
			}
		}
	}
}

func (mgr *Manager) StartService(hash string) error {
	// todo: check running instances by hash...
	service, err := loadServiceFromDB(hash)
	if err != nil {
		return err
	}
	utils.PrintDebug("manager: service %s loaded from db", service.Hash)
	m, err := mgr.startService(service)
	if err != nil {
		log.Printf("failed to start service (%s)", err)
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

func (mgr *Manager) StartTunnel(hash string) (*tunnels.Tunnel, error) {
	// todo: check running instances by hash...
	tunnel, err := loadTunnelFromDB(hash)
	if err != nil {
		return nil, err
	}

	if utils.IsDevMode() {
		log.Printf("manager: tunnel %s loaded from db", tunnel.Hash)
	}

	err = mgr.startTunnel(tunnel)
	return nil, err
}

func (mgr *Manager) StopTunnel(hash string) error {
	return mgr.stopTunnel(hash)
}
