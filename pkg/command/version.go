package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// VersionCommand is a struct that represent version command.
type VersionCommand struct {
}

// Show version
func (v *VersionCommand) Show(c *cli.Context) error {
	fmt.Println("0.0.6")
	return nil
}
