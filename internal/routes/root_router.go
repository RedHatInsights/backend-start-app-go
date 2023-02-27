package routes

import "github.com/go-chi/chi/v5"

func RootRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Mount(PathPrefix(), apiRouter())

	return router
}
