package main

import (
	"log"
	"os"

	"github.com/apicat/apicat/v2"
	"github.com/apicat/apicat/v2/backend/migrations"
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
				Value:   ".env",
				Usage:   "Load configuration from `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "make:migration",
				Usage: "Automated DB schema updates",
				Action: func(cCtx *cli.Context) error {
					migrations.Generate(cCtx.Args().First())
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			app := apicat.NewApp(ctx.String("config"))
			return app.Run()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
