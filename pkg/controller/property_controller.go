package controller

import (
	"net/http"
	"os"

	"github.com/TsuyoshiUshio/volley/pkg/helper"
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
)

func UpdateJMeterConfig(c *gin.Context) {
	var jMeterProperty model.JMeterProperty
	if err := c.ShouldBindJSON(&jMeterProperty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := os.Stat(model.JMeterPropertyFile); err == nil {
		err = os.Remove(model.JMeterPropertyFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	err := jMeterProperty.GenerateModifiedProperty(helper.GetJMeterPropertyFilePath(), model.JMeterPropertyFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jMeterProperty)
	return
}
