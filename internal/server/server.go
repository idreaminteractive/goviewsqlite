// Main app, takes in db, context, doncif and other items for testing and main app serving
package server

import (
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/idreaminteractive/goviewsqlite/internal/sqlite"
)

type Server struct {
	DB         *sqlite.LocalDB
	Mux        *chi.Mux
	logger     *httplog.Logger
	HttpServer *http.Server
}

func NewServer(logger *httplog.Logger, database *sqlite.LocalDB) *Server {
	if !database.Ready() {
		log.Fatalf("Database is not ready")
	}

	r := SetupMux(logger)

	return &Server{
		DB:     database,
		Mux:    r,
		logger: logger,
		// init http after
	}
}

type MuxSetupStruct struct {
	Mux *chi.Mux
}

func (m *Server) Initialize(mux *chi.Mux, logger *httplog.Logger, portStr string) error {

	// web.AddRootMiddlewares(m.Mux,  m.logger, )

	// attach + create all handlers
	addRoutes(m.Mux, m.logger, m.DB.Connection())

	port, err := strconv.Atoi(portStr)
	if err != nil {
		m.logger.Error("Could not parse port")
		return err
	}
	// init the server (but not run it!)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m.Mux,
	}

	m.HttpServer = httpServer
	return nil

}

// only run during a normal run site... not in tests!
func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	if s.HttpServer == nil {
		return errors.New("HttpServer is nil")
	}

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("err: %v", err)
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// if s.config.IsLocalDev() {
	// 	hrUrl, err := common.HotReloadUrl(s.config)
	// 	if err != nil {
	// 		log.Fatalf("Could not get site url for hot reload")
	// 	}
	// 	err = goreload.SendReloadSignal(hrUrl)

	// 	if err != nil {
	// 		log.Fatalf("Failed to trigger hot reload: %v", err)
	// 	}

	// }
	s.logger.Info("GoSQLite Viewer started!")
	// url, err := common.SiteUrl(s.config)

	// if err == nil {
	// 	fmt.Printf("URL: %s\n", url)
	// }

	var wg sync.WaitGroup

	// for each goroutine to wait, we add one...
	wg.Add(1)
	go func() {
		defer wg.Done()
		// wait for our done signal (ctrl+c or otherwise)
		<-ctx.Done()

		// graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.HttpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		// close our db too.
		if s.DB != nil {
			if err := s.DB.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing database: %s\n", err)
			}
		}
	}()
	//
	wg.Wait()
	return nil
}

// take in our interfaces + setup all the details here.
func SetupMux(logger *httplog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// add our logger to chi + to server struct
	if logger != nil { // for tests
		r.Use(httplog.RequestLogger(logger))
	}

	// setup our enforcer here
	// clean path (users//2 -> users/2)
	r.Use(middleware.CleanPath)

	// remove trailing slashes
	r.Use(middleware.StripSlashes)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))

	validation.ErrorTag = "form"

	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
		// Debug:            true,
		MaxAge: 300,
	}))

	return r
}
