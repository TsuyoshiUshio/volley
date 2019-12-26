package model

// Config is a struct for storing config information of JMeter
type Job struct {
	ID string `json:"id" binding:"required"`
}

type JobResponse struct {
	ConfigID string `json:"config_id" binding:"required"`
	JobID    string `json:"job_id" binding:"required"`
}

type JobRequest struct {
	ConfigID      string `json:"config_id" binding:"required"`
	IsDistributed bool   `json:"is_distributed" binding:"required"`
}
