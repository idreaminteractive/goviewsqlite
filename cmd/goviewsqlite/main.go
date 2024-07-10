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
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "seed",
				Usage: "Locally sets up a db for usage for testing @ the db location",
				Action: func(cCtx *cli.Context) error {
					// cmd := exec.Command("cat", "./data/seed.sql", "|", "sqlite3", "/data/sqlite.db")
					// cmd := exec.Command("sqlite3", "-cmd", ".import ./data/seed.sql")
					// stdoutStderr, err := cmd.CombinedOutput()
					// if err != nil {
					// 	fmt.Fprintf(cCtx.App.ErrWriter, "exec failed  %s: \n%s\n", err, string(stdoutStderr))
					// } else {
					// 	fmt.Fprintf(cCtx.App.Writer, "Result: %s", string(stdoutStderr))
					// }

					// return nil
					return commands.RunSeed(cCtx.App.Writer, url, "./data/seed.sql")
				},
			},
		},
	}

	app.Run(os.Args)

}
