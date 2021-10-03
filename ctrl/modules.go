package ctrl

import (
	"github.com/vikingo-project/vsat/modules"

	"github.com/gin-gonic/gin"
)

func hmodules(c *gin.Context) {
	avaliableModules := modules.GetAvaliableModules()
	var modules []map[string]interface{}
	for _, m := range avaliableModules {
		modules = append(modules, m.GetInfo())
	}

	c.JSON(200, gin.H{"status": "ok", "modules": modules})
}
