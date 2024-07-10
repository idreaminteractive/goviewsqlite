// Command interface for running the hot reload server
package main

import (
	"fmt"
	"os"

	"github.com/idreaminteractive/goviewsqlite/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	var url string
	var port string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dbpath",
				Value:       "/data/sqlite.db",
				Usage:       "Location of the database",
				Destination: &url,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       "5555",
				Usage:       "Port to listen on",
				Destination: &port,
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Printf("Starting up - %s, %s\n", url, port)
			return commands.RunServer(cCtx.Context, cCtx.App.Writer, url, port)
		},
		Commands: []*cli.Command{
			{
				Name:  "seed",
				Usage: "Locally sets up a db for usage for testing @ the db location",
				Action: func(cCtx *cli.Context) error {
					return commands.RunSeed(cCtx.App.Writer, url, "./data/seed.sql")
				},
			},
		},
	}

	app.Run(os.Args)

}
