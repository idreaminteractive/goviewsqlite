package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/idreaminteractive/goviewsqlite/internal/logger"
	"github.com/idreaminteractive/goviewsqlite/internal/server"
	"github.com/idreaminteractive/goviewsqlite/internal/sqlite"
	_ "modernc.org/sqlite"
)

func RunServer(ctx context.Context, w io.Writer, dbPath, port string) error {

	fmt.Fprintln(w, "Starting up")

	logger := logger.NewLogger(w)
	db := sqlite.NewDB(dbPath)

	if err := db.Open(); err != nil {
		logger.Error("Could not open db", "error", err)
		return err
	}

	s := server.NewServer(logger, db)
	s.Initialize(s.Mux, logger, port)

	// ok, let's run it.
	if err := s.Run(ctx); err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return err
	}
	return nil

}
