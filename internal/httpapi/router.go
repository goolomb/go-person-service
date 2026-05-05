package httpapi

import "github.com/go-chi/chi/v5"

func NewRouter(personService PersonSaver) chi.Router {
	router := chi.NewRouter()
	handler := NewHandler(personService)

	router.Get("/health", HealthHandler)
	router.Post("/save", handler.SaveHandler)
	router.Get("/{id}", handler.GetByIDHandler)

	return router
}
