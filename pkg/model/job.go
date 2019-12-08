package model

// Config is a struct for storing config information of JMeter
type Job struct {
	ID string `json:"id" binding:"required"` 
}