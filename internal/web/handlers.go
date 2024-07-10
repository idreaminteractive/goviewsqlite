package web

import (
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/idreaminteractive/goviewsqlite/internal/views/components"
)

// Handlers for base routes (health, root, etc)
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	htmx.NewResponse().RenderTempl(r.Context(), w, components.BasePageComponent(components.BasePageData{
		Title:   "Not Found!",
		Content: components.NotFoundComponent(),
	}))
}

// base routes
func HandleAnyHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func HandleRootGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		component := components.RootComponent(components.RootComponentProps{})
		htmx.NewResponse().RenderTempl(r.Context(), w, components.BasePageComponent(components.BasePageData{
			Title:   "Welcome to Ignite",
			Content: component,
		}))

	}
}
