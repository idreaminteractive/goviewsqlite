// Command interface for running the hot reload server
package main

import (
	"fmt"
	"os"

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
			return nil
		},
		// Commands: []*cli.Command{
		// 	{
		// 		Name:  "server",
		// 		Usage: "Runs hot reload server on a specific port",
		// 		Action: func(cCtx *cli.Context) error {

		// 			return commands.Serve(cCtx, url)
		// 		},
		// 	},
		// },
	}

	app.Run(os.Args)

}
