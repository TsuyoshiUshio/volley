package controller

import (
	"net/http"
	"github.com/TsuyoshiUshio/volley/pkg/provider"
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Start(c *gin.Context) {
	var config model.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job_id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	runContext := provider.NewRunContext(config.ID, job_id.String())
	p := provider.NewAzureProvider()

	go func() {
		p.Run(runContext)
	}()

	c.JSON(http.StatusCreated, gin.H{
		"job_id": job_id.String(),
		"config_id": config.ID,
	})

	return

}
