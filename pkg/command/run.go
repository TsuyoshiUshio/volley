package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	return nil
	// TODO you can implate mode that wait until the execution is finished by polling status api.
}
