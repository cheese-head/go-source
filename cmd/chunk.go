package cmd

import (
	"context"
	"encoding/json"
	"go-source/pkg/parser"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func Chunk() *cli.Command {
	return &cli.Command{
		Name:  "chunk",
		Usage: "chunk files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "output file path",
			},

			&cli.StringFlag{
				Name:     "files",
				Usage:    "glob pattern for files that need to be chunked",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "include",
				Usage: "list of file extensions to include",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "set log level to debug",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			globPattern := ctx.String("files")
			includes := ctx.StringSlice("include")
			output := ctx.String("output")
			debug := ctx.IsSet("debug")
			if debug {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}

			file := os.Stdout
			var fileError error
			if output != "" {
				file, fileError = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if fileError != nil {
					return fileError
				}
			}
			readFilesWithGlobPattern(globPattern, includes, file)
			return nil
		},
	}
}

func processPath(path string, includes []string, writer io.Writer) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			for _, include := range includes {
				if filepath.Ext(path) == include {
					chunker, err := parser.DetectLanguageFromFile(path)
					if err != nil {
						return err
					}
					result, err := parser.ParseFile(context.Background(), chunker, path)
					if err != nil {
						return err
					}
					data, err := json.Marshal(result)
					if err != nil {
						return err
					}
					_, err = writer.Write(data)
					if err != nil {
						return err
					}

					_, err = writer.Write([]byte("\n"))
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func readFilesWithGlobPattern(globPattern string, includes []string, writer io.Writer) error {
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		err := processPath(match, includes, writer)
		if err != nil {
			return err
		}
	}

	return nil
}
