package commands

import (
	"context"
	"fmt"
	"io"

	"os"

	"github.com/idreaminteractive/goviewsqlite/internal/logger"
	_ "modernc.org/sqlite"
)

func RunServer(ctx context.Context, w io.Writer, dbPath, port string) error {

	fmt.Println("Starting up")

	// init our stuff we'll be using

	logger := logger.NewLogger(w)

	// create the server

	// ok, let's run it.
	if err := s.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
