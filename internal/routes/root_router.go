package routes

import "github.com/go-chi/chi/v5"

func RootRouter() *chi.Mux {
	router := chi.NewRouter()

	apiR := apiRouter()
	router.Mount(pathVersionedPrefix("v1"), apiR)
	router.Mount(pathVersionedPrefix("v1.0"), apiR)

	return router
}
