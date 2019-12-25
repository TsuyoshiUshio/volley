package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/TsuyoshiUshio/volley/pkg/helper"
	"github.com/TsuyoshiUshio/volley/pkg/model"
)

type Provider interface {
}

type AzureProvider struct {
}
type RunContext struct {
	ConfigID      string
	JobID         string
	IsDistributed bool
}

type StatusCheckContext struct {
	JobID string
}

func NewAzureProvider() *AzureProvider {
	return &AzureProvider{}
}

func NewStatusCheckContext(jobID string) *StatusCheckContext {
	return &StatusCheckContext{
		JobID: jobID,
	}
}

func NewRunContext(configID string, jobID string, isDistributed bool) *RunContext {
	return &RunContext{
		ConfigID:      configID,
		JobID:         jobID,
		IsDistributed: isDistributed,
	}
}
func (p *AzureProvider) StatusCheck(context *StatusCheckContext) (*model.Status, error) {
	statusFilePath := filepath.Join(".", model.JobDir, context.JobID, model.StatusFileName)
	statusFile, err := ioutil.ReadFile(statusFilePath)
	if err != nil {
		return nil, err
	}
	var status model.Status
	err = json.Unmarshal(statusFile, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
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

	// Transfer csv files to the Slave not if necessary
	if context.IsDistributed {
		err = transferCSVFilesToSlaveIfJMeterPropertyJSONExists(configDirPath)
		if err != nil {
			return err
		}
	}

	cmd, err := buildCommand(configFilePath, logPath, reportPath, context.IsDistributed)
	if err != nil {
		return err
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

// CreateAsset returns Zipped file
func (p *AzureProvider) CreateAsset(jobID string) ([]byte, error) {
	dir, err := ioutil.TempDir("", "volley")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	sourcePath := filepath.Join(".", model.JobDir, jobID)
	zipFilePath := filepath.Join(dir, jobID+".zip")
	helper.Zip(sourcePath, zipFilePath)
	return ioutil.ReadFile(zipFilePath)
}

func buildCommand(configFilePath, logPath, reportPath string, isDistributed bool) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	jmeterProperty, err := model.NewJMeterProperty()
	if err != nil && isDistributed {
		return nil, fmt.Errorf("Can not find Remote_IPs configration. Use property API to configure the settings. error: %v", err)
	}
	if isDistributed && jmeterProperty.CommaSeparatedRemoteHostIPs() == "" {
		return nil, fmt.Errorf("Remote_IPs configuration is empty. Use property API to configure the settings")
	}

	if runtime.GOOS == "windows" {
		commandPath, err := exec.LookPath("jmeter.bat")
		if err != nil {
			return nil, err
		}
		if isDistributed {
			cmd = exec.Command("cmd", "/C", commandPath, "-n", "-t", configFilePath, "-l", logPath, "-e", "-o", reportPath, "-R", jmeterProperty.CommaSeparatedRemoteHostIPs())
		} else {
			cmd = exec.Command("cmd", "/C", commandPath, "-n", "-t", configFilePath, "-l", logPath, "-e", "-o", reportPath)
		}
	} else {
		if isDistributed {
			cmd = exec.Command("jmeter", "-n", "-t", configFilePath, "-l", logPath, "-e", "-o", reportPath, "-R", jmeterProperty.CommaSeparatedRemoteHostIPs())
		} else {
			cmd = exec.Command("jmeter", "-n", "-t", configFilePath, "-l", logPath, "-e", "-o", reportPath)
		}
	}
	return cmd, nil
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

func transferCSVFilesToSlaveIfJMeterPropertyJSONExists(configDirPath string) error {
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	fieldName := "file"
	matchCsv := regexp.MustCompile(`.*\.csv`)
	err := filepath.Walk(configDirPath, func(path string, info os.FileInfo, err error) error {
		if matchCsv.MatchString(path) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			fw, err := mw.CreateFormFile(fieldName, path)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, file)
			if err != nil {
				return err
			}
			file.Close()
		}
		return nil
	})
	if err != nil {
		return err
	}
	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	err = postMultipartTotheSlaves(contentType, body)
	if err != nil {
		return err
	}
	return nil
}

func postMultipartTotheSlaves(contentType string, body []byte) error {
	jmeterConfigJSONPath := filepath.Join(".", model.JMeterPropertyDir, model.JMeterPropertyJSON)
	if _, err := os.Stat(jmeterConfigJSONPath); err != nil {
		// JMeterConfig JSON exists
		file, err := ioutil.ReadFile(jmeterConfigJSONPath)
		if err != nil {
			return err
		}
		var jmeterProperty model.JMeterProperty
		err = json.Unmarshal(file, &jmeterProperty)
		if err != nil {
			return err
		}
		for _, ipAddress := range jmeterProperty.RemoteHostIPs {
			err = postToSlave(ipAddress, contentType, body)
		}
		return nil
	} else {
		// JMeterConfig JSON does not exists
		return nil
	}
}

func postToSlave(iPAddress, contentType string, body []byte) error {
	url := fmt.Sprintf("http://%s:%s/csv", iPAddress, model.SlaveDefaultPort)
	buf := &bytes.Buffer{}
	buf.Write(body)
	resp, err := http.Post(url, contentType, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(responseBody))
	return nil
}
