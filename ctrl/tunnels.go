package ctrl

import (
	"log"
	"time"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func httpTunnels(c *gin.Context) {
	var tuns []models.Tunnel
	dq := db.GetConnection().Model(&models.Tunnel{})

	err := dq.Find(&tuns).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatal(err)
	}

	mgr := getManager(c)
	for i, t := range tuns {
		live := mgr.Tunnels.Exists(t.Hash)
		tuns[i].Connected = live
		if live {
			tuns[i].PublicAddr = mgr.Tunnels.GetPublicAddr(t.Hash)
		}
	}
	c.JSON(200, gin.H{"status": "ok", "tunnels": tuns})
}

func httpCreateTunnel(c *gin.Context) {
	type p struct {
		DstHost   string `json:"dstHost" binding:"required"`
		DstPort   int    `json:"dstPort" binding:"required"`
		DstTLS    bool   `json:"dstTLS"`
		Type      string `json:"type" binding:"required"`
		AutoStart bool   `json:"autoStart"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	err := db.GetConnection().Save(&models.Tunnel{
		Hash:      utils.EasyHash(true),
		Type:      params.Type,
		DstHost:   params.DstHost,
		DstPort:   params.DstPort,
		DstTLS:    params.DstTLS,
		Autostart: params.AutoStart,
	}).Error

	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func httpUpdateTunnel(c *gin.Context) {
	type p struct {
		Hash      string `json:"hash" binding:"required"`
		DstHost   string `json:"dstHost" binding:"required"`
		DstPort   int    `json:"dstPort" binding:"required"`
		DstTLS    bool   `json:"dstTLS"`
		Type      string `json:"type" binding:"required"`
		AutoStart bool   `json:"autoStart"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	err := db.GetConnection().Model(&models.Tunnel{}).Where(&models.Tunnel{Hash: params.Hash}).Updates(&models.Tunnel{
		DstHost: params.DstHost, DstPort: params.DstPort, DstTLS: params.DstTLS, Autostart: params.AutoStart}).Error
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func httpRemoveTunnel(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// restart service
	mgr := getManager(c)
	mgr.StopTunnel(params.Hash)

	// waiting for 1s, service should be stopped
	time.Sleep(time.Second)
	db.GetConnection().Delete(&models.Tunnel{}, &models.Tunnel{Hash: params.Hash})
	c.JSON(200, gin.H{"status": "ok"})
}

func httpStartTunnel(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	mgr := getManager(c)
	_, err := mgr.StartTunnel(params.Hash)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// waiting for public addr...
	time.Sleep(time.Second)
	c.JSON(200, gin.H{"status": "ok"})
}

func httpStopTunnel(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}
	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	mgr := getManager(c)
	mgr.StopTunnel(params.Hash)
	c.JSON(200, gin.H{"status": "ok"})
}
