package controller

import (
	"net/http"

	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
)

// OverrideJMeterProperty write JMeterProperty as a configration file
func OverrideJMeterProperty(c *gin.Context) {
	var jMeterProperty model.JMeterProperty
	if err := c.ShouldBindJSON(&jMeterProperty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// save the JMeterProperty
	if err := jMeterProperty.WriteFile(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	return
}
