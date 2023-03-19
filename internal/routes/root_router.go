package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func RootRouter() *chi.Mux {
	router := chi.NewRouter()

	apiR := apiRouter()

	// Set Content-Type to JSON for chi renderer. Warning: Non-chi routes
	// MUST set Content-Type header on their own!
	apiR.Use(render.SetContentType(render.ContentTypeJSON))

	router.Mount(pathVersionedPrefix("v1"), apiR)
	router.Mount(pathVersionedPrefix("v1.0"), apiR)

	return router
}
