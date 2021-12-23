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
	filterDates := c.QueryArray("dates[]")

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

	if len(filterDates) > 0 {
		if len(filterDates) == 2 {
			tq.Where("date BETWEEN ? AND ?", filterDates[0], filterDates[1])
			dq.Where("date BETWEEN ? AND ?", filterDates[0], filterDates[1])
		}
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

// httpEvents is a HTTP handler for getting events from the interaction
func httpEvents(c *gin.Context) {
	hash := c.Param("hash")
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize := 10 // events per page
	offset := (page - 1) * pageSize

	var total int64
	ec := db.GetConnection().Model(&models.FullEvent{}).Where(&models.FullEvent{Session: hash})
	ec.Count(&total)

	var events []models.FullEvent
	ed := db.GetConnection().Model(&models.Session{})
	err := ed.Model(&models.FullEvent{}).Where(&models.FullEvent{Session: hash}).Order("date ASC").Limit(pageSize).Offset(offset).Find(&events).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(200, gin.H{"status": "error", "error": err.Error()})
			return
		}
	}

	// mark session as visited
	db.GetConnection().Model(&models.Session{}).Where("hash = ?", hash).Update("visited", true)
	c.JSON(200, gin.H{"status": "ok", "events": events, "total": total})
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
