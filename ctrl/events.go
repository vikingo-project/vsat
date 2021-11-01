package ctrl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func hSessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("size"))
	interHash := c.Query("hash")

	filterService := c.QueryArray("service[]")
	filterClientIP := strings.TrimSpace(c.Query("client_ip"))
	filterLocalAddr := strings.TrimSpace(c.Query("local_addr"))
	filterDescription := strings.TrimSpace(c.Query("description"))

	if page < 1 {
		page = 1
	}
	if pageSize < 15 {
		pageSize = 15
	}

	offset := (page - 1) * pageSize
	var total int64
	tq := db.GetConnection().Model(&models.Session{})
	dq := db.GetConnection().Model(&models.Session{})
	var sessions []models.Session
	if interHash != "" {
		tq.Where("hash == ?", interHash)
		dq.Where("hash == ?", interHash)
	}

	if len(filterService) > 0 {
		tq.Where("service in (?)", filterService)
		dq.Where("service in (?)", filterService)
	}

	if filterClientIP != "" {
		tq.Where("client_ip LIKE ?", fmt.Sprintf("%%%s%%", filterClientIP))
		dq.Where("client_ip LIKE ?", fmt.Sprintf("%%%s%%", filterClientIP))
	}

	if filterLocalAddr != "" {
		tq.Where("local_addr LIKE ?", fmt.Sprintf("%%%s%%", filterLocalAddr))
		dq.Where("local_addr LIKE ?", fmt.Sprintf("%%%s%%", filterLocalAddr))
	}

	if filterDescription != "" {
		tq.Where("description LIKE ?", fmt.Sprintf("%%%s%%", filterDescription))
		dq.Where("description LIKE ?", fmt.Sprintf("%%%s%%", filterDescription))
	}

	tq.Count(&total)
	err := dq.Order("date DESC").Limit(pageSize).Offset(offset).Find(&sessions).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(200, gin.H{"status": "error", "error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{"status": "ok", "sessions": sessions, "total": total})
}

func hEvents(c *gin.Context) {
	hash := c.Param("hash")
	var events []models.FullEvent
	err := db.GetConnection().Model(&models.FullEvent{}).Where(&models.FullEvent{Session: hash}).Find(&events).Error
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// set session visited
	db.GetConnection().Model(&models.Session{}).Where("hash = ?", hash).Update("visited", true)
	c.JSON(200, gin.H{"status": "ok", "events": events})
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
