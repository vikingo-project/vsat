package ctrl

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/models"
)

func hFiles(c *gin.Context) {
	fileHash := c.Query("hash")
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize := 50
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var total int64
	cq := db.GetConnection().Model(&models.File{})
	cq.Count(&total)

	dq := db.GetConnection().Model(&models.File{}).Select("hash,date,file_name,size,content_type,interaction_hash,service_hash").Order("id() desc").Limit(pageSize).Offset(offset)
	if fileName != "" {
		dq.Where("file_name LIKE ?", strings.TrimSpace(fileName))
	}
	if fileType != "" {
		dq.Where("content_type == ?", fileType)
	}
	if fileHash != "" {
		dq.Where("hash == ?", fileHash)
	}

	var files []models.File
	err := dq.Find(&files)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "total": total, "files": files})
}

func hTypes(c *gin.Context) {
	var files []models.File
	db.GetConnection().Model(&models.File{}).Select("distinct content_type").Find(&files)

	var types []string
	for _, f := range files {
		types = append(types, f.ContentType)
	}
	c.JSON(200, gin.H{"status": "ok", "types": types})
}

func hDownloadFile(c *gin.Context) {
	hash := c.Param("hash")
	var file models.File
	db.GetConnection().Select("data,content_type").Where(&models.File{Hash: hash}).Find(&file)
	c.Data(http.StatusOK, file.ContentType, file.Data)
	c.Abort()
}

func hUploadFiles(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	filename := filepath.Base(file.Filename)
	f, _ := file.Open()
	buff := make([]byte, file.Size)
	f.Read(buff)
	f.Close()

	fileStruct := files.PrepareFile(filename, buff)
	err = db.GetConnection().Begin().Save(&fileStruct)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
}

func hRemoveFile(c *gin.Context) {
	type p struct {
		Hash string `json:"hash" binding:"required"`
	}
	var params p
	if err := c.ShouldBindJSON(&params); err == nil {
		err := db.GetConnection().Begin().Delete(&models.File{}, &models.File{Hash: params.Hash})
		if err != nil {
			c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": "ok"})
		}
	} else {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
	}
}
