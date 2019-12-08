package model

// Config is a struct for storing config information of JMeter
type Config struct {
	ID string `json:"id" binding:"required"` 
}

const (
	// ConfigDir is a directory that contains config files
	ConfigDir = "config"
	// JobDir is a directory that contains job reports
	JobDir = "job"
)