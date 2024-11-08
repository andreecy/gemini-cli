package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gmn",
		Usage: "AI Powered CLI",
		Commands: []*cli.Command{
			{
				Name:  "ask",
				Usage: "Ask any question to the AI",
				Action: func(cCtx *cli.Context) error {
					prompt := strings.Join(cCtx.Args().Slice(), " ")
					generateContent(prompt)
					return nil
				},
			},
			{
				Name:  "howto",
				Usage: "General-purpose how to do something in terminal",
				Action: func(cCtx *cli.Context) error {
					prompt := strings.Join(cCtx.Args().Slice(), " ")
					howto(prompt)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
