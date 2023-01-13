package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
	"gorm.io/gorm"
)

func (a *APIC) Sessions(params string) (*RecordsContainer, error) {
	q, _ := url.ParseQuery(params)
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("size"))
	interHash := q.Get("hash")

	filterService := utils.QueryArray(q, "service[]")
	filterClientIP := strings.TrimSpace(q.Get("client_ip"))
	filterLocalAddr := strings.TrimSpace(q.Get("local_addr"))
	filterDescription := strings.TrimSpace(q.Get("description"))
	filterDates := utils.QueryArray(q, "dates[]")

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
			return &RecordsContainer{}, err
		}
	}
	return &RecordsContainer{Total: total, Records: sessions}, nil
}

func (a *APIC) SessionEvents(params string) (*RecordsContainer, error) {
	p, _ := url.ParseQuery(params)
	page, _ := strconv.Atoi(p.Get("page"))
	pageSize, _ := strconv.Atoi(p.Get("size"))
	hash := p.Get("hash")
	offset := (page - 1) * pageSize

	if page < 1 {
		page = 1
	}

	if pageSize < 15 {
		pageSize = 15
	}

	var total int64
	ec := db.GetConnection().Model(&models.FullEvent{}).Where(&models.FullEvent{Session: hash})
	ec.Count(&total)

	var events []models.FullEvent
	ed := db.GetConnection().Model(&models.Session{})
	err := ed.Model(&models.FullEvent{}).Where(&models.FullEvent{Session: hash}).Order("date ASC").Limit(pageSize).Offset(offset).Find(&events).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return &RecordsContainer{}, err
		}
	}
	// mark session as visited
	db.GetConnection().Model(&models.Session{}).Where("hash = ?", hash).Update("visited", true)
	return &RecordsContainer{Records: events, Total: total}, err
}
