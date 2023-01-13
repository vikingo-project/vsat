package ctrl

import (
	"github.com/vikingo-project/vsat/api"

	"github.com/gin-gonic/gin"
)

func httpAbout(c *gin.Context) {
	res := api.Instance.About()
	c.JSON(200, res)
}
