package main

import (
	"go-source/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name: "go-source",
		Commands: []*cli.Command{
			cmd.Chunk(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
