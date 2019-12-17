package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/urfave/cli/v2"
)

type RunCommand struct {
}

func (s *RunCommand) Run(c *cli.Context) error {
	configID := c.String("config-id")
	masterIP := c.String("master")
	port := c.String("port")
	if port == "" {
		port = "38080"
	}
	requestBody, err := json.Marshal(&model.Config{
		ID: configID,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(masterIP+":"+port+"/job", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if c.Bool("wait") {
		timeoutDuration := c.Int("timeout")

		var jobRequest model.JobRequest
		json.Unmarshal(body, &jobRequest)
		ch := make(chan string, 1)
		go func() {
			status := waitForCompletion(masterIP, port, jobRequest.JobID)
			ch <- status
		}()

		select {
		case status := <-ch:
			if status == model.StatusFailed {
				os.Exit(1) // If the status is failed, exit with 1.
			}
			fmt.Println("Done.")
		case <-time.After(time.Duration(timeoutDuration) * time.Minute):
			fmt.Printf("timeout! : %d minutes.\n", timeoutDuration)
		}
	}

	return nil
	// TODO you can implate mode that wait until the execution is finished by polling status api.
}

func waitForCompletion(masterIP, port, jobID string) string {
	defaultPollingInterval := 5
	fmt.Println("Waiting for Job completion...")
	var status model.Status
	start := time.Now()

	for {
		time.Sleep(time.Duration(defaultPollingInterval) * time.Second)
		requestURI := masterIP + ":" + port + "/job/" + jobID
		resp, err := http.Get(requestURI) // Ignore the error for retrying.
		elapsed := time.Since(start)
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Polling status for JobID: %s Failed reading body, error: %v \n", jobID, err)
			}
			json.Unmarshal(body, &status)
			if status.Status == model.StatusCompleted || status.Status == model.StatusFailed {
				fmt.Printf("Finish execution JobID: %s Status: %s in %s second\n", jobID, status.Status, elapsed)
				return status.Status
			}
			fmt.Printf("Polling status for JobID: %s Status: %s at %s second ...\n", jobID, status.Status, elapsed)
		} else {
			fmt.Printf("Polling status for JobID: %s Faile request, error: %v \n", jobID, err)
		}
	}
}
