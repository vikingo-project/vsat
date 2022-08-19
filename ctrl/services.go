package ctrl

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vikingo-project/vsat/api"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/modules"
	"github.com/vikingo-project/vsat/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func httpServices(c *gin.Context) {
	res, err := api.Instance.Services(c.Request.URL.RawQuery)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Total": res.Total, "Records": res.Records})
}

func httpCreateService(c *gin.Context) {
	type p struct {
		ListenIP    string      `json:"listenIP" binding:"required"`
		ListenPort  int         `json:"listenPort" binding:"required"`
		ServiceName string      `json:"serviceName" binding:"required"`
		AutoStart   bool        `json:"autoStart"`
		ModuleName  string      `json:"moduleName" binding:"required"`
		Settings    interface{} `json:"moduleSettings" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// todo: validate settings...
	settings, _ := json.Marshal(params.Settings)

	module := modules.GetModuleByName(params.ModuleName)
	baseProto := strings.Join(module.GetInfo()["base_proto"].([]string), "/")

	err := db.GetConnection().Save(&models.Service{
		Hash:        utils.EasyHash(false),
		ServiceName: params.ServiceName,
		ModuleName:  params.ModuleName,
		ListenIP:    params.ListenIP,
		ListenPort:  params.ListenPort,
		Autostart:   params.AutoStart,
		Settings:    string(settings),
		BaseProto:   baseProto,
	}).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func httpUpdateService(c *gin.Context) {
	type p struct {
		Hash        string      `json:"hash" binding:"required"`
		ServiceName string      `json:"serviceName" binding:"required"`
		ListenIP    string      `json:"listenIP" binding:"required"`
		ListenPort  int         `json:"listenPort" binding:"required"`
		AutoStart   bool        `json:"autoStart"`
		Settings    interface{} `json:"moduleSettings" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	settings, _ := json.Marshal(params.Settings)
	err := db.GetConnection().Model(&models.Service{}).Where(&models.Service{Hash: params.Hash}).Updates(&models.Service{ServiceName: params.ServiceName,
		ListenIP: params.ListenIP, ListenPort: params.ListenPort, Autostart: params.AutoStart, Settings: string(settings)}).Error
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func httpRemoveService(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// restart service

	manager.StopService(params.Hash)

	// waiting for 1s, service should be stopped
	time.Sleep(time.Second)
	db.GetConnection().Delete(&models.Service{}, &models.Service{Hash: params.Hash})
	c.JSON(200, gin.H{"status": "ok"})
}

func httpStartService(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	err := manager.StartService(params.Hash)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func httpStopService(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}

	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	manager.StopService(params.Hash)
}
