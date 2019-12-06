package provider

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TsuyoshiUshio/volley/pkg/model"
)

type Provider interface {
}

type AzureProvider struct {
}
type RunContext struct {
	ConfigID    string
	JobID       string
	JmxFileName string
}

func (p AzureProvider) Run(context *RunContext) error {
	jobPath := filepath.Join(".", model.JobDir, context.JobID)
	if _, err := os.Stat(jobPath); os.IsNotExist(err) {
		os.Mkdir(jobPath, os.ModeDir)
	}

	configDirPath := filepath.Join(".", model.ConfigDir, context.ConfigID)
	configFilePath := filepath.Join(configDirPath, context.JmxFileName)
	statusPath := filepath.Join(".", model.JobDir, context.JobID)
	status := model.Status{
		Status: model.StatusRunning,
	}
	status.Write(statusPath)
	cmd := exec.Command("jmeter", "-n", "-t", configFilePath, "-l", "stress.log", "-e", "-o", "report")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		status = model.Status{
			Status: model.StatusFailed,
		}
		status.Write(statusPath)
		return err
	}
	status = model.Status {
		Status: model.StatusCompleted,
	}
	status.Write(statusPath)
	return nil
}
