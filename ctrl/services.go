package ctrl

import (
	"net/http"
	"time"

	"github.com/vikingo-project/vsat/api"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/manager"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"

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

func httpCreateUpdateService(c *gin.Context) {
	var service models.WebService
	c.Bind(&service)
	var (
		res string
		err error
	)
	if service.Hash == "" {
		utils.PrintDebug("create a new service")
		res, err = api.Instance.CreateService(&service)
	} else {
		utils.PrintDebug("update service")
		res, err = api.Instance.UpdateService(&service)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "Hash": res})
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

func httpToggleService(c *gin.Context) {
	var params models.ChangeServiceState
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	err := api.Instance.ToggleService(&params)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})

}
