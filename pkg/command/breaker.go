package command

import (
	"log"
	"os"

	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/urfave/cli/v2"
)

// BreakerCommand is a sub command for breaking build if an execution result doesn't meet a success criteria.
type BreakerCommand struct {
}

// Validate method validate if it an execution result meets a success criteria that is specified
func (b *BreakerCommand) Validate(c *cli.Context) error {
	logFilePath := c.String("log-file")
	configPath := c.String("config")

	criteria, err := model.NewAverageTimeAndErrorOnRPSSuccessCriteria(configPath)
	if err != nil {
		return err
	}
	validation, err := criteria.Validate(logFilePath)
	if err != nil {
		return err
	}

	if !validation {
		log.Println("Validation failed.")
		os.Exit(1) // finish with status code 1.
	} else {
		log.Println("Validation succeed. ")
	}
	return nil
}
