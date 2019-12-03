package command

import (
	"github.com/urfave/cli/v2"
	"fmt"
)

type ProvisionCommand struct {	
}

func (s *ProvisionCommand) Provision(c *cli.Context) error {
	fmt.Println("Provisioning JMeter Environment Master: 1, Slave:", c.Int("slave"))
	return nil
}