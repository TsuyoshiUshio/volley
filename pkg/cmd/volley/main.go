package main

import (
	"os"
	"github.com/urfave/cli/v2"
	"log"
	"github.com/TsuyoshiUshio/volley/pkg/command"
)

func main() {
	app := &cli.App{
		Name: "volley", 
		Usage: "Manage JMeter cluster",
		Version: "0.0.1",
		Commands: []*cli.Command {
			{
				Name: "provision",
				Aliases: []string{"p"},
				Usage: "Provision JMeter cluster",
				Action: (&command.ProvisionCommand{}).Provision,
				Flags: []cli.Flag{
					&cli.IntFlag {
						Name: "slave",
						Aliases: []string{"s"},
						Value: 1,
						Usage: "Specify the number of slaves of JMeter cluster",
					},
				}, 
			},
		},		
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}