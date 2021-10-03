package ctrl

import (
	"log"
	"strconv"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"

	"github.com/gin-gonic/gin"
)

func hSessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("size"))
	interHash := c.Query("hash")

	filterService := c.QueryArray("service[]")
	filterClientIP := c.Query("client_ip")
	filterLocalAddr := c.Query("local_addr")
	filterDescription := c.Query("description")

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
		tq.Where("client_ip LIKE ?", filterClientIP)
		dq.Where("client_ip LIKE ?", filterClientIP)
	}

	if filterLocalAddr != "" {
		tq.Where("local_addr LIKE ?", filterLocalAddr)
		dq.Where("local_addr LIKE ?", filterLocalAddr)
	}

	if filterDescription != "" {
		tq.Where("description LIKE ?", filterDescription)
		dq.Where("description LIKE ?", filterDescription)
	}

	tq.Count(&total)
	err := dq.Order("id() desc").Limit(pageSize).Offset(offset).Find(&sessions)
	if err != nil {
		log.Print(err)
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "sessions": sessions, "total": total})
}

func hEvents(c *gin.Context) {
	hash := c.Param("hash")

	var events []models.FullEvent
	err := db.GetConnection().Model(&models.FullEvent{}).Where(&models.FullEvent{Session: hash}).Find(&events)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// set session visited
	db.GetConnection().Model(&models.Session{}).Where("hash = ?", hash).Update("visited", true)
	c.JSON(200, gin.H{"status": "ok", "events": events})
}
