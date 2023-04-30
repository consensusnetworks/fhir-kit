package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

const (
	defaultPort = "9090"
)

type Config struct {
	Port        string
	Verbose     bool
	InContainer bool
	CGOEnabled  bool
}

func Run() error {
	app := cli.NewApp()

	app.Name = "fhir-kit"
	app.Description = "Simple-to-spin-up FHIR development kit"
	app.Version = "1.0.0"
	app.BashComplete = cli.DefaultAppComplete
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose output",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "inContainer",
			Usage: "running in a container",
			Value: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		var err error
		// err := godotenv.Load()

		// if err != nil {
		// 	return err
		// }

		config := Config{
			Port:        os.Getenv("PORT"),
			Verbose:     c.Bool("verbose"),
			InContainer: c.Bool("inContainer"),
			CGOEnabled:  os.Getenv("CGO_ENABLED") == "false",
		}

		if config.Port == "" {
			config.Port = defaultPort
		}

		server := NewServer(config)

		err = server.Start()

		if err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		return err
	}

	return nil
}
