package controller

import (
	"log"
	"net/http"
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
		fileName := filepath.Base(file.Filename)
		log.Println(fileName)
		dist := filepath.Join(csvPath, fileName)
		err := c.SaveUploadedFile(file, dist)
		if err != nil {
			log.Printf("Csv File upload error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
