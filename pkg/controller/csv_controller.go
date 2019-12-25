package controller

import (
	"log"
	"os"
	"path/filepath"

	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
)

// UploadCSV enalbe us to get the multipart file that is assumed csv file.
// Save to to the csv dir on the current directory
func UploadCSV(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	csvPath := filepath.Join(".", model.CsvDir)
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		err = os.MkdirAll(csvPath, os.ModePerm)
	}
	for _, file := range files {
		log.Println(file.Filename)
		dist := filepath.Join(csvPath, file.Filename)
		c.SaveUploadedFile(file, dist)
	}
}
