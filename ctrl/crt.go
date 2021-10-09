package ctrl

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

func hCerts(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize := 15
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var total int64
	db.GetConnection().Model(&models.Crt{}).Count(&total)
	utils.PrintDebug("hCerts: total certificates %d", total)

	var certs []models.Crt
	err := db.GetConnection().Model(&models.Crt{}).Select("name").Order("id desc").Limit(pageSize).Offset(offset).Find(&certs)
	if err != nil {
		log.Print(err)
	}
	c.JSON(200, gin.H{"status": "ok", "total": total, "certs": certs})
}
