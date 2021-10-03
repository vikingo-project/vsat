package ctrl

import (
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/shared"
	"github.com/vikingo-project/vsat/utils"

	"github.com/gin-gonic/gin"
)

func about(c *gin.Context) {
	c.JSON(200, gin.H{"version": shared.Version, "build": shared.BuildHash})
}

// sql is a handler that is helps to debug SQL
func sql(c *gin.Context) {
	if utils.IsDevMode() {
		q := c.Query("q")
		rows, err := db.SQLQuery(q)
		c.JSON(200, gin.H{"rows": rows, "err": err})
		return
	}
	c.JSON(200, gin.H{})
}
