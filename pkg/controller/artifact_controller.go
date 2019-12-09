package controller

import (
	"bytes"
	"net/http"

	"github.com/TsuyoshiUshio/volley/pkg/provider"
	"github.com/gin-gonic/gin"
)

func Download(c *gin.Context) {
	// Create a zip from the target directry and zip it as a tempfile.
	p := provider.NewAzureProvider()
	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can not file job_id as a parameter. /artifact/:job_id required."})
	}
	assetFile, err := p.CreateAsset(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the stream from the tempfile
	b := bytes.NewBuffer(assetFile)
	contentLength := b.Len()
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + jobID + `.zip"`,
	}
	c.DataFromReader(http.StatusOK, int64(contentLength), "application/zip", b, extraHeaders)
}
