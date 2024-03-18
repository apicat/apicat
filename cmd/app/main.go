package main

import (
	apicatcloud "apicat-cloud"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "apicat",
		Usage: "run apicat app",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.yaml",
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(ctx *cli.Context) error {
			app := apicatcloud.NewApp(ctx.String("config"))
			return app.Run()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
