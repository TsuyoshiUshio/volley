package controller

import (
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/TsuyoshiUshio/volley/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func Start(c *gin.Context) {
	var request model.JobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job_id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	runContext := provider.NewRunContext(request.ConfigID, job_id.String(), request.IsDistributed)
	p := provider.NewAzureProvider()

	go func() {
		p.Run(runContext)
	}()

	c.JSON(http.StatusCreated, gin.H{
		"job_id":    job_id.String(),
		"config_id": request.ConfigID,
	})

	return

}

func StatusCheck(c *gin.Context) {

	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not file job_id as a parameter. /job/:job_id required."})
		return
	}
	p := provider.NewAzureProvider()
	statusCheckContext := provider.NewStatusCheckContext(jobID)
	status, err := p.StatusCheck(statusCheckContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
