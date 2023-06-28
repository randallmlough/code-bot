package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "take-this",
		Usage: "takes an input and prompt and does something with it using GPT",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "the source input",
				FilePath: "testdata/input.ts",
			},
			&cli.StringFlag{
				Name:     "prompt",
				Aliases:  []string{"p"},
				Usage:    "the prompt for what to do with the input",
				FilePath: "testdata/prompt.txt",
			},
		},
		Action: Action,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
