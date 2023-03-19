package routes

import (
	"consoledot-go-template/api"
	"consoledot-go-template/internal/logging"
	"consoledot-go-template/internal/services"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// PathPrefix determines the entry point for all public incoming requests.
// The server prefixes all API paths with this prefix.
func PathPrefix() string {
	return "/api/template"
}

func pathVersionedPrefix(version string) string {
	return fmt.Sprintf("%s/%s", PathPrefix(), version)
}

func apiRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(logging.NewMiddleware(log.Logger))
	mountSpec(router)
	mountAPI(router)
	return router
}

func mountSpec(router *chi.Mux) {
	router.Get("/openapi.json", api.ServeOpenAPISpec)
}

func mountAPI(router *chi.Mux) {
	router.Route("/hellos", func(r chi.Router) {
		r.Get("/", services.ListHellos)
	})
}
