package controller

func UpdateJMeterConfig(c *gin.Context) {
	var jMeterProperty model.JMeterProperty
	if err := c.ShouldBindJSON(&jMeterProperty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Implement the method that enable up to update the JMeter Property File
	
}