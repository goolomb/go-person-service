package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goolomb/go-person-service/internal/service"
)

const (
	errorAlreadyExists = "already_exists"
	errorInternal      = "internal_error"
	errorInvalidJSON   = "invalid_json"
	errorNotFound      = "not_found"
	errorValidation    = "validation_error"
	jsonContentType    = "application/json"
	contentTypeHeader  = "Content-Type"
	locationHeader     = "Location"
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
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	var request SavePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   errorInvalidJSON,
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
		writeJSON(w, http.StatusBadRequest, ValidationErrorResponse{
			Error:   errorValidation,
			Message: "Request body contains validation errors.",
			Fields:  validationErrors,
		})
		return
	}
	if errors.Is(err, service.ErrAlreadyExists) {
		writeJSON(w, http.StatusConflict, ErrorResponse{
			Error:   errorAlreadyExists,
			Message: "Person already exists.",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   errorInternal,
			Message: "An unexpected error occurred.",
		})
		return
	}

	w.Header().Set(locationHeader, "/"+person.ExternalID)
	writeJSON(w, http.StatusCreated, personResponseFromService(person))
}

func (h Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	person, validationErrors, err := h.personService.GetPersonByExternalID(chi.URLParam(r, "id"))
	if errors.Is(err, service.ErrValidation) {
		writeJSON(w, http.StatusBadRequest, ValidationErrorResponse{
			Error:   errorValidation,
			Message: "Request path contains validation errors.",
			Fields:  validationErrors,
		})
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   errorNotFound,
			Message: "Person was not found.",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   errorInternal,
			Message: "An unexpected error occurred.",
		})
		return
	}

	writeJSON(w, http.StatusOK, personResponseFromService(person))
}

func personResponseFromService(person service.Person) PersonResponse {
	return PersonResponse{
		ExternalID:  person.ExternalID,
		Name:        person.Name,
		Email:       person.Email,
		DateOfBirth: person.DateOfBirth,
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
