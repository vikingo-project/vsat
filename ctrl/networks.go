package ctrl

import (
	"github.com/vikingo-project/vsat/utils"

	"github.com/gin-gonic/gin"
)

func networks(c *gin.Context) {
	networks, _ := utils.GetNetworks()
	c.JSON(200, gin.H{"status": "ok", "networks": networks})
}
