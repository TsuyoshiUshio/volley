package model

import (
	"os"
	"encoding/json"
	"path/filepath"
)

type Status struct {
	Status string `json:"status"`
}

const (
	StatusRunning   = "running"
	StatusCompleted = "complated"
	StatusFailed    = "failed"
)

const(
	StatusFileName = "status.json"
)

func (s Status) Write(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	statusPath := filepath.Join(path, StatusFileName)
	statusJson, _ := json.Marshal(s)
	f, err := os.Create(statusPath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(statusJson))
	if err != nil {
		return err
	}
	return nil
}
