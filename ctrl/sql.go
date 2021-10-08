package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/utils"
)

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
