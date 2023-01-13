package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
)

func (a *APIC) Files(params string) (*RecordsContainer, error) {
	p, _ := url.ParseQuery(params)
	page, _ := strconv.Atoi(p.Get("page"))
	pageSize, _ := strconv.Atoi(p.Get("size"))
	fileHash := p.Get("hash")
	fileName := strings.TrimSpace(p.Get("file_name"))
	fileType := p.Get("file_type")

	var (
		total  int64
		offset = (page - 1) * pageSize
	)

	cq := db.GetConnection().Model(&models.File{})
	cq.Count(&total)

	dq := db.GetConnection().Model(&models.File{}).Select("hash,date,file_name,size,content_type,interaction_hash,service_hash").Order("date DESC").Limit(pageSize).Offset(offset)
	if fileName != "" {
		dq.Where("file_name LIKE ?", fmt.Sprintf("%%%s%%", fileName))
	}
	if fileType != "" {
		dq.Where("content_type == ?", fileType)
	}
	if fileHash != "" {
		dq.Where("hash == ?", fileHash)
	}

	var files []models.File
	err := dq.Find(&files).Error
	return &RecordsContainer{Total: total, Records: files}, err
}

func (a *APIC) FileTypes() (*RecordsContainer, error) {
	var files []models.File
	err := db.GetConnection().Model(&models.File{}).Select("distinct content_type").Find(&files).Error
	var types []string
	for _, f := range files {
		types = append(types, f.ContentType)
	}
	return &RecordsContainer{Records: types}, err
}
