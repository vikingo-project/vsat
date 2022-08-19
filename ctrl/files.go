package ctrl

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/vikingo-project/vsat/api"
	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/files"
	"github.com/vikingo-project/vsat/models"
)

func httpFiles(c *gin.Context) {
	res, err := api.Instance.Files(c.Request.URL.RawQuery)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "Total": res.Total, "Records": res.Records})
}

func httpFileTypes(c *gin.Context) {
	res, err := api.Instance.FileTypes()
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "Records": res.Records})
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
	err = db.GetConnection().Save(&fileStruct).Error
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
		err := db.GetConnection().Delete(&models.File{}, &models.File{Hash: params.Hash}).Error
		if err != nil {
			c.JSON(200, gin.H{"status": "error", "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": "ok"})
		}
	} else {
		c.JSON(200, gin.H{"status": "error", "error": err.Error()})
	}
}
