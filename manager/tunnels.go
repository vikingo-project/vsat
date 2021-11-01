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
	if utils.IsDevMode() {
		log.Println("Start tunnel", tunnel.Hash)
	}

	t := &tunnels.Tunnel{
		Hash:           tunnel.Hash,
		Type:           tunnel.Type,
		Destination:    fmt.Sprintf("%s:%d", tunnel.DstHost, tunnel.DstPort),
		DestinationTLS: tunnel.DstTLS,
	}

	var (
		errChan = make(chan error)
	)

	err := t.Start(errChan)
	if err != nil {
		log.Println("failed to start tunnel ", err)
		return err
	}

	go func() {
		defer log.Println("exit from waiter routine")
		<-errChan
		utils.PrintDebug("err chan", err)
		mgr.Tunnels.Remove(tunnel.Hash)
	}()

	// mark tunnel as live
	mgr.Tunnels.Add(tunnel.Hash, t)
	return nil
}

func (mgr *Manager) stopTunnel(hash string) error {
	if utils.IsDevMode() {
		log.Println("Stop tunnel", hash)
	}
	t := mgr.Tunnels.Get(hash)
	if t != nil {
		t.Stop()
		mgr.Tunnels.Remove(hash)
	}
	return nil
}
