package manager

import (
	"errors"
	"fmt"
	"log"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/tunnels"
	"github.com/vikingo-project/vsat/utils"
)

func loadTunnelsFromDB() ([]models.Tunnel, error) {
	var tunnels []models.Tunnel
	err := db.GetConnection().Find(&tunnels).Error
	return tunnels, err
}

func loadTunnelFromDB(hash string) (models.Tunnel, error) {
	var (
		tunnel = models.Tunnel{Hash: hash}
		count  int64
	)

	db.GetConnection().Model(&tunnel).Where(&tunnel).Count(&count)
	if count != 1 {
		return tunnel, errors.New("tunnel not found")
	}

	err := db.GetConnection().Find(&tunnel, &tunnel).Error
	return tunnel, err
}

func (mgr *Manager) startTunnel(tunnel models.Tunnel) error {
	utils.PrintDebug("Start tunnel %s", tunnel.Hash)

	t := &tunnels.Tunnel{
		Hash:           tunnel.Hash,
		Type:           tunnel.Type,
		Destination:    fmt.Sprintf("%s:%d", tunnel.DstHost, tunnel.DstPort),
		DestinationTLS: tunnel.DstTLS,
	}

	var errChan = make(chan error)

	err := t.Start(errChan)
	if err != nil {
		log.Println("failed to start tunnel ", err)
		return err
	}

	go func(hash string) {
		defer utils.PrintDebug("remove tunnel %s", hash)
		<-errChan
		mgr.Tunnels.Remove(hash)
	}(tunnel.Hash)

	mgr.Tunnels.Add(tunnel.Hash, t)
	return nil
}

func (mgr *Manager) stopTunnel(hash string) error {
	utils.PrintDebug("Stop tunnel %s", hash)
	t := mgr.Tunnels.Get(hash)
	if t != nil {
		t.Stop()
		mgr.Tunnels.Remove(hash)
	}
	return nil
}
