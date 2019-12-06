package model

import (
	"os"
)

type Status struct {
	Status string `json:"status"`
}



const (
	StatusRunning = "running"
	StatusCompleted = "complated"
	StatusFailed = "failed"
)

func (s Status) Write(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, os.ModeDir)
	}

	configDirPath := jilepath.Join(".", model.ConfigDir, json.ID)

	statusPath = filePath.Join(filePath, "status.json")
	statusJson, _ := json.Marshall(s)
	f, err := os.Create(statusPath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(statusJson)
	if err != nil {
		return err
	}
	return nil
}