package routes

import (
	"consoledot-go-template/internal/services"

	"github.com/go-chi/chi/v5"
)

// PathPrefix determines the entry point for all public incoming requests.
// The server prefixes all API paths with this prefix.
func PathPrefix() string {
	return "/api/template"
}

func apiRouter() *chi.Mux {
	router := chi.NewRouter()
	mountAPI(router)
	return router
}

func mountAPI(router *chi.Mux) {
	router.Route("/hellos", func(r chi.Router) {
		r.Get("/", services.ListHellos)
	})
}
