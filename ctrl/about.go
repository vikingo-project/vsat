package ctrl

import (
	"github.com/vikingo-project/vsat/shared"

	"github.com/gin-gonic/gin"
)

func about(c *gin.Context) {
	c.JSON(200, gin.H{"version": shared.Version, "build": shared.BuildHash})
}
