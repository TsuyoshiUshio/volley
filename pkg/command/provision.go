package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"time"
)

type ProvisionCommand struct {
}

func (s *ProvisionCommand) Provision(c *cli.Context) error {
	fmt.Printf("Provisioning JMeter Cluster. Name: %s Master: 1, Slave: %d ...", c.String("cluster-name"), c.Int("slave"))
	time.Sleep(5 * time.Second)
	return nil
}
