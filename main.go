package main

import (
	"github.com/therealkevinard/adr-er/commands"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "adr-er",
		Usage: "a friendly little thing for managing architectural decision records",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:        "create",
				Aliases:     []string{"new"},
				Usage:       "create a new adr document",
				Description: "new is used to create a brand-spankin-new adr document",
				Action:      func(c *cli.Context) error { return new(commands.Create).Action(c) },
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("oops. app exited with an error:", err)
	}

}
