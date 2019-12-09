package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNewConfig(c *gin.Context) {
	// Receive multiple files
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	configID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	configPath := filepath.Join(".", model.ConfigDir, configID.String())
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configPath, os.ModePerm)
	}
	for _, file := range files {
		log.Println(file.Filename)
		dist := filepath.Join(configPath, file.Filename)
		c.SaveUploadedFile(file, dist)
	}
	c.JSON(http.StatusOK, &model.Config{
		ID: configID.String(),
	})

}
