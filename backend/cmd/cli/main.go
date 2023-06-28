package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "code-bot",
		Usage: "a helpful coding assistant",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "input",
			},
			&cli.StringFlag{
				Name:    "prompt",
				Aliases: []string{"p"},
				Usage:   "the prompt for what to do with the input",
			},
		},
		Action: Action,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
