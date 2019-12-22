package main

import (
	"github.com/TsuyoshiUshio/volley/pkg/command"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "volley",
		Usage:   "Load Testing Tool with JMeter",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:    "provision",
				Aliases: []string{"p"},
				Usage:   "Provision JMeter cluster",
				Action:  (&command.ProvisionCommand{}).Provision,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cluster-name",
						Aliases: []string{"c"},
						Usage:   "Specify Cluster Name. Should be uniq.",
					},
					&cli.IntFlag{
						Name:    "slave",
						Aliases: []string{"s"},
						Value:   1,
						Usage:   "Specify the number of slaves of JMeter cluster",
					},
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "API Server for uploading/receiving files",
				Action:  (&command.ServerCommand{}).Start,
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run JMeter",
				Action:  (&command.RunCommand{}).Run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config-id",
						Aliases: []string{"c"},
						Usage:   "Specify config-id that is created by config command.",
					},
					&cli.StringFlag{
						Name:    "master",
						Aliases: []string{"m"},
						Usage:   "Specify master ip address or domain name.",
					},
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "38080",
						Usage:   "Specify master port. 38080 by default",
					},
					&cli.BoolFlag{
						Name:    "wait",
						Aliases: []string{"w"},
						Usage:   "Make this subcommand wait for completion",
					},
					&cli.IntFlag{
						Name:    "timeout",
						Aliases: []string{"t"},
						Value:   30,
						Usage:   "Specify the default timeout in minutes if you use --wait (-w) flag",
					},
					&cli.StringFlag{
						Name:    "output-type",
						Aliases: []string{"o"},
						Value:   "stdout",
						Usage:   "Specify the how to output the job_id. Possible value is 'stdout', 'file', 'both', if you choose file or both, it will output as file. The file name will respect outpus-filename flag",
					},
					&cli.StringFlag{
						Name:    "output-filename",
						Aliases: []string{"of"},
						Value:   "job.json",
						Usage:   "Specify the output filename when you specify --output-type flag",
					},
				},
			}, {
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Upload jmx, csv files to the server. Return value is config-id.",
				Action:  (&command.ConfigCommand{}).Upload,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "directory",
						Aliases: []string{"d"},
						Usage:   "Specify directory that contains jmx and csv files that you want to upload",
					},
					&cli.StringFlag{
						Name:    "master",
						Aliases: []string{"m"},
						Usage:   "Specify master ip address or domain name.",
					},
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "38080",
						Usage:   "Specify master port. 38080 by default",
					},
				},
			},
			{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "fetch log of the JMeter run.",
				Action:  (&command.LogCommand{}).Download,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "job-id",
						Aliases: []string{"j"},
						Usage:   "Specify job-id that run sub command returns.",
					},
					&cli.StringFlag{
						Name:    "master",
						Aliases: []string{"m"},
						Usage:   "Specify master ip address or domain name.",
					},
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "38080",
						Usage:   "Specify master port. 38080 by default",
					},
				},
			},
			{
				Name:    "breaker",
				Aliases: []string{"b"},
				Usage:   "Build ",
				Action:  (&command.BreakerCommand{}).Validate,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "log-file",
						Aliases: []string{"l"},
						Usage:   "File path for JMeter execution log file.",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "success_criteria.json",
						Usage:   "Config file path of success_criteria",
					},
				},
			},
			{
				Name:    "destroy",
				Aliases: []string{"d"},
				Usage:   "Destroy the JMeter Cluster",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cluster-name",
						Aliases: []string{"c"},
						Usage:   "Specify Cluster Name. Should be uniq.",
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
