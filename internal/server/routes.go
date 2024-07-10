package server

import (
	"database/sql"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/idreaminteractive/goviewsqlite/internal/web"
)

// pass in concrete handlers that have the access in the respective services
func addRoutes(
	r *chi.Mux,
	logger *httplog.Logger,
	db *sql.DB,
) {

	// static files come from static folder. :O
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.NotFound(web.HandleNotFound)

	r.Get("/", web.HandleRootGet())

}
