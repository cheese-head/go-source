package cmd

import (
	"go-source/pkg/server"

	"github.com/urfave/cli/v2"
)

func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "start a http server",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "set log level to debug",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			return server.Run()
		},
	}
}
