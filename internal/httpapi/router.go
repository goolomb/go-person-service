package httpapi

import "github.com/go-chi/chi/v5"

func NewRouter() chi.Router {
	router := chi.NewRouter()

	router.Get("/health", HealthHandler)
	router.Post("/save", SaveHandler)
	router.Get("/{id}", GetByIDHandler)

	return router
}
