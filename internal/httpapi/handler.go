package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goolomb/go-person-service/internal/service"
)

type PersonSaver interface {
	SavePerson(input service.SavePersonInput) (service.Person, map[string]string, error)
	GetPersonByExternalID(id string) (service.Person, map[string]string, error)
}

type Handler struct {
	personService PersonSaver
}

func NewHandler(personService PersonSaver) Handler {
	return Handler{personService: personService}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	var request SavePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "invalid_json",
			Message: "Request body must be a valid JSON object.",
		})
		return
	}

	person, validationErrors, err := h.personService.SavePerson(service.SavePersonInput{
		ExternalID:  request.ExternalID,
		Name:        request.Name,
		Email:       request.Email,
		DateOfBirth: request.DateOfBirth,
	})
	if errors.Is(err, service.ErrValidation) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ValidationErrorResponse{
			Error:   "validation_error",
			Message: "Request body contains validation errors.",
			Fields:  validationErrors,
		})
		return
	}
	if errors.Is(err, service.ErrAlreadyExists) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "already_exists",
			Message: "Person already exists.",
		})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "An unexpected error occurred.",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/"+person.ExternalID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(personResponseFromService(person))
}

func (h Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	person, validationErrors, err := h.personService.GetPersonByExternalID(chi.URLParam(r, "id"))
	if errors.Is(err, service.ErrValidation) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ValidationErrorResponse{
			Error:   "validation_error",
			Message: "Request path contains validation errors.",
			Fields:  validationErrors,
		})
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "not_found",
			Message: "Person was not found.",
		})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "internal_error",
			Message: "An unexpected error occurred.",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(personResponseFromService(person))
}

func personResponseFromService(person service.Person) PersonResponse {
	return PersonResponse{
		ExternalID:  person.ExternalID,
		Name:        person.Name,
		Email:       person.Email,
		DateOfBirth: person.DateOfBirth,
	}
}
