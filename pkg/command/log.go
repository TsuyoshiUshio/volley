package command

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/TsuyoshiUshio/volley/pkg/helper"
	"github.com/urfave/cli/v2"
)

type LogCommand struct {
}

func (l *LogCommand) Download(c *cli.Context) error {
	jobID := c.String("job-id")
	masterIP := c.String("master")
	port := c.String("port")

	resp, err := http.Get(masterIP + ":" + port + "/asset/" + jobID)
	if err != err {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create(jobID + ".zip")

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = helper.UnZip(jobID+".zip", ".")
	if err != nil {
		return err
	}
	defer func() {
		err := out.Close()
		if err != nil {
			log.Println("close:", err)
		}
		err = os.Remove(jobID + ".zip")
		if err != nil {
			log.Println("remove:", err)
		}
	}()
	return nil
}
