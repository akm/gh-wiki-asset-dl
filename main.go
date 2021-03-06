package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "target-url",
				Value: "https://user-images.githubusercontent.com",
			},
			&cli.IntFlag{
				Name:    "download-buffer-size",
				Aliases: []string{"b"},
				Value:   1024 * 1024,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "version",
				Usage:  "Show version",
				Action: showVersion,
			},
		},
		Action: execute,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func showVersion(ctx *cli.Context) error {
	fmt.Fprintf(os.Stdout, "%s\n", VERSION)
	return nil
}

func execute(ctx *cli.Context) error {
	paths := ctx.Args()
	if paths.Len() < 1 {
		if err := cli.ShowAppHelp(ctx); err != nil {
			return err
		}
		return fmt.Errorf("No path given")
	}

	ff := NewFilter([]string{".md"})
	rep := NewReplacer(
		NewScanner(ctx.String("target-url")),
		NewDownloader(ctx.Int("download-buffer-size")),
	)

	for _, path := range paths.Slice() {
		if err := ff.Glob(path, rep.Do); err != nil {
			return err
		}
	}

	return nil
}
