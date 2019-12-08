package provider

import (
	"fmt"
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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

func NewAzureProvider() *AzureProvider {
	return &AzureProvider{}
}

func NewRunContext(configID string, jobID string) *RunContext {
	return &RunContext{
		ConfigID: configID,
		JobID:    jobID,
	}
}

func (p *AzureProvider) Run(context *RunContext) error {
	dir, _ := os.Getwd()
	fmt.Println(dir)
	jobPath := filepath.Join(".", model.JobDir, context.JobID)
	if _, err := os.Stat(jobPath); os.IsNotExist(err) {
		err = os.MkdirAll(jobPath, os.ModePerm)
	}

	configDirPath := filepath.Join(".", model.ConfigDir, context.ConfigID)
	configFileName, err := findJmxFile(configDirPath)
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(configDirPath, configFileName)
	status := model.Status{
		Status: model.StatusRunning,
	}
	status.Write(jobPath)
	// TODO Parameter should be flexible.
	logPath := filepath.Join(jobPath, "stress.log")
	reportPath := filepath.Join(jobPath, "report")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		commandPath, err := exec.LookPath("jmeter.bat")
		if err != nil {
			return err
		}
		cmd = exec.Command("cmd", "/C", commandPath, "-n", "-t", configFilePath, "-l", logPath, "-e", "-o", reportPath)

	} else {
		cmd = exec.Command("jmeter", "-n", "-t", configFilePath, "-l", logPath, "-e", "-o", reportPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		status = model.Status{
			Status: model.StatusFailed,
		}
		status.Write(jobPath)
		return err
	}
	status = model.Status{
		Status: model.StatusCompleted,
	}
	status.Write(jobPath)
	return nil
}

func findJmxFile(directory string) (string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`.*\.jmx`)
	for _, f := range files {
		if r.MatchString(f.Name()) {
			return f.Name(), nil
		}
	}
	return "", fmt.Errorf("Can not file jmx file on the config directory: %s", directory)
}
