// Main app, takes in db, context, doncif and other items for testing and main app serving
package server

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"idreaminteractive/ignite/internal/config"
	errs "idreaminteractive/ignite/internal/errors"
	"idreaminteractive/ignite/internal/features/authz"
	"idreaminteractive/ignite/internal/features/common"
	"idreaminteractive/ignite/internal/features/game"
	"idreaminteractive/ignite/internal/features/jobs"
	"idreaminteractive/ignite/internal/features/unity"

	"idreaminteractive/ignite/internal/features/organization"
	"idreaminteractive/ignite/internal/features/turso"
	"idreaminteractive/ignite/internal/features/user"
	ignite "idreaminteractive/ignite/internal/types"
	"idreaminteractive/ignite/internal/web"

	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/idreaminteractive/goreload"
)

type IgniteDB interface {
	Connection() *sql.DB
	Open() error
	Close() error
	Ready() bool
}

type MuxConfigStruct struct {
	SessionSecret string
	DisableCSRF   bool
	Logger        *httplog.Logger
	Env           string
}

type Server struct {
	DB                IgniteDB
	Mux               *chi.Mux
	config            *config.EnvConfig
	logger            *httplog.Logger
	HttpServer        *http.Server
	RunningServices   []RunningService
	CloseableServices []CloseableService
}

// Services that are created and passed in by the calling system to NewServer
type Services struct {
	// Non-concrete or variable services should meet an interface defined in here and the consuming services
	TenantService ignite.TenantService

	// for concrete services, require a pointer directly to them
	OrgService   *organization.OrganizationService
	AuthzService *authz.AuthzService
	IdpService   ignite.IDPService
	// this may be moved to non-concrete
	SessionService ignite.SessionService

	UserService *user.UserService
	// may move to non-concrete
	TursoService *turso.TursoService
	GameService  *game.GameService
	UnityService *unity.UnityService

	JobService *jobs.JobService
}

// used to check if services are nil - they cannot be.
func anyFieldIsNil(s interface{}) bool {
	v := reflect.ValueOf(s)

	// Handle pointer to struct
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			return true
		}
	}
	return false
}

// this gets run when this module is loaded. it's hacky,
// but gob needs these types setup before anything else
// so, it works.
func init() {
	gob.Register(ignite.SessionPayload{})
	gob.Register(ignite.FlashMessage{})
}

func NewServer(config *config.EnvConfig, logger *httplog.Logger, database IgniteDB, disableCSRF bool) *Server {
	if !database.Ready() {
		log.Fatalf("Database is not ready")
	}

	r := SetupMux(MuxConfigStruct{
		SessionSecret: config.SessionSecret,
		Logger:        logger,
		Env:           config.DopplerConfig,
		DisableCSRF:   disableCSRF,
	})

	return &Server{
		DB:     database,
		Mux:    r,
		config: config,
		logger: logger,
		// init http after
	}
}

func (m *Server) Initialize(svs *Services) error {
	if anyFieldIsNil(svs) {
		m.logger.Error("Field in service struct was nil")
		return errs.Errorf(errs.INTERNAL, "Field in service struct was nil")
	}
	web.AddRootMiddlewares(m.Mux, m.config, m.logger, svs.SessionService, svs.OrgService, svs.AuthzService, svs.GameService, svs.TenantService)

	// attach + create all handlers
	addRoutes(m.Mux, m.config, m.logger, svs)

	port, err := strconv.Atoi(m.config.Port)
	if err != nil {
		m.logger.Error("Could not parse port")
		return err
	}
	// init the server (but not run it!)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m.Mux,
	}

	// are there services that need to be shutdown?
	m.RunningServices = append(m.RunningServices, svs.JobService)

	m.HttpServer = httpServer
	return nil

}

type RunningService interface {
	Run(ctx context.Context) error
}

type CloseableService interface {
	Close() error
}

// only run during a normal run site... not in tests!
func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	if s.HttpServer == nil {
		return errs.Errorf(errs.INTERNAL, "HttpServer is nil, cannot run")
	}
	s.logger.Info("ignite started!")
	url, err := common.SiteUrl(s.config)

	if err == nil {
		fmt.Printf("URL: %s\n", url)
	}

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("err: %v", err)
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	if s.config.IsLocalDev() {
		hrUrl, err := common.HotReloadUrl(s.config)
		if err != nil {
			log.Fatalf("Could not get site url for hot reload")
		}
		err = goreload.SendReloadSignal(hrUrl)

		if err != nil {
			log.Fatalf("Failed to trigger hot reload: %v", err)
		}

	}

	// run our services
	for _, runnable := range s.RunningServices {
		// we run the service w/ the ctx, so if we shutdown, it'll capture it
		if err := runnable.Run(ctx); err != nil {
			log.Fatalf("Failed to start runnable service: %v", err)
		}
	}

	var wg sync.WaitGroup

	// for each goroutine to wait, we add one...
	wg.Add(1)
	go func() {
		defer wg.Done()
		// wait for our done signal (ctrl+c or otherwise)
		<-ctx.Done()

		// close all closables
		for _, closable := range s.CloseableServices {
			if err := closable.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing service: %v\n", err)
			}
		}

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
func SetupMux(
	config MuxConfigStruct,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// add our logger to chi + to server struct
	if config.Logger != nil { // for tests
		r.Use(httplog.RequestLogger(config.Logger))
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
