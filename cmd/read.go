package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-source/pkg/parser"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func Read() *cli.Command {
	return &cli.Command{
		Name:  "read",
		Usage: "read content from chunks",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "set log level to debug",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			scanner := bufio.NewScanner(os.Stdin)
			writer := os.Stdout

			if ctx.IsSet("debug") {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
			for scanner.Scan() {
				text := scanner.Text()
				result := &parser.Result{}
				err := json.Unmarshal([]byte(text), result)
				if err != nil {
					slog.Error("unable to unmarshal input", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
					return err
				}
				for _, v := range result.Chunks {
					writer.Write(v.Content)
					writer.Write([]byte("\n"))
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				os.Exit(1)
			}
			return nil
		},
	}
}
