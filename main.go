package main

import (
	"go-source/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:        "go-source",
		Description: "A tool for parsing and chunking source code",
		Commands: []*cli.Command{
			cmd.Chunk(),
			cmd.Read(),
			cmd.Server(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
