package files

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/vikingo-project/vsat/db"
	"github.com/vikingo-project/vsat/models"
	"github.com/vikingo-project/vsat/utils"
)

func GetFileByHash(hash string, getBody bool) models.File {
	var (
		columns = "hash,file_name,size,content_type"
		file    models.File
	)
	if getBody {
		columns += ",data"
	}
	err := db.GetConnection().Select(columns).Where(&models.File{Hash: hash}).Find(&file)
	if err != nil {
		log.Print("failed to get file info", err)
	}
	return file
}

func PrepareFile(name string, data []byte) models.File {
	hash := utils.EasyHash(true)
	fileHeader := make([]byte, 512)
	copy(fileHeader, data)
	contentType := http.DetectContentType(fileHeader)
	return models.File{
		Date:        time.Now(),
		Hash:        hash,
		FileName:    name,
		ContentType: contentType,
		Size:        int64(len(data)),
		Data:        data,
	}
}

func GetFilenameFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
