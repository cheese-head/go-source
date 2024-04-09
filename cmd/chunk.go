package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"go-source/pkg/parser"
	"io/fs"
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
				Name:     "files",
				Usage:    "glob pattern for files that need to be chunked",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "include",
				Usage: "list of file extensions to include",
			},
		},
		Action: func(ctx *cli.Context) error {
			globPattern := ctx.String("files")
			includes := ctx.StringSlice("include")
			readFilesWithGlobPattern(globPattern, includes)
			return nil
		},
	}
}

func processPath(path string, includes []string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			for _, include := range includes {
				if filepath.Ext(path) == include {
					result, err := parser.ParseFile(context.Background(), path)
					if err != nil {
						return err
					}
					data, _ := json.Marshal(result)
					fmt.Println(string(data))
				}
			}
		}
		return nil
	})
}

func readFilesWithGlobPattern(globPattern string, excludes []string) error {
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		fmt.Println(err)

		return err
	}

	for _, match := range matches {
		err := processPath(match, excludes)
		if err != nil {
			fmt.Println(err)

			return err
		}
	}

	return nil
}
