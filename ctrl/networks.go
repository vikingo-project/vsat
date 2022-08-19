package ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vikingo-project/vsat/api"
)

func httpNetworks(c *gin.Context) {
	res, err := api.Instance.Networks()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "Total": 0, "Records": make(map[string]string, 0)})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Total": res.Total, "Records": res.Records})
}
