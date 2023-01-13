package ctrl

import (
	"net/http"

	"github.com/vikingo-project/vsat/api"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"

	"github.com/gin-gonic/gin"
)

func httpSessions(c *gin.Context) {
	res, err := api.Instance.Sessions(c.Request.URL.RawQuery)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Total": res.Total, "Records": res.Records})
}

// httpEvents is a HTTP handler for getting events from the interaction
func httpEvents(c *gin.Context) {
	res, err := api.Instance.SessionEvents(c.Request.URL.RawQuery)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Total": res.Total, "Records": res.Records})
}

func hRemoveSession(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}
	var params p
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	db.GetConnection().Where("hash = ?", params.Hash).Delete(&models.Session{})
	c.JSON(200, gin.H{"status": "ok"})
}
