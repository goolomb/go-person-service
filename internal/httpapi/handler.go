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

type PersonService interface {
	SavePerson(input service.SavePersonInput) (service.Person, map[string]string, error)
	GetPersonByExternalID(id string) (service.Person, map[string]string, error)
}

type Handler struct {
	personService PersonService
}

func NewHandler(personService PersonService) Handler {
	return Handler{personService: personService}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	var request SavePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, errorInvalidJSON, "Request body must be a valid JSON object.")
		return
	}

	person, validationErrors, err := h.personService.SavePerson(savePersonInputFromRequest(request))
	if errors.Is(err, service.ErrValidation) {
		writeValidationError(w, "Request body contains validation errors.", validationErrors)
		return
	}
	if errors.Is(err, service.ErrAlreadyExists) {
		writeError(w, http.StatusConflict, errorAlreadyExists, "Person already exists.")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, errorInternal, "An unexpected error occurred.")
		return
	}

	w.Header().Set(locationHeader, "/"+person.ExternalID)
	writeJSON(w, http.StatusCreated, personResponseFromService(person))
}

func (h Handler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	person, validationErrors, err := h.personService.GetPersonByExternalID(chi.URLParam(r, "id"))
	if errors.Is(err, service.ErrValidation) {
		writeValidationError(w, "Request path contains validation errors.", validationErrors)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, errorNotFound, "Person was not found.")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, errorInternal, "An unexpected error occurred.")
		return
	}

	writeJSON(w, http.StatusOK, personResponseFromService(person))
}

func savePersonInputFromRequest(request SavePersonRequest) service.SavePersonInput {
	return service.SavePersonInput{
		ExternalID:  request.ExternalID,
		Name:        request.Name,
		Email:       request.Email,
		DateOfBirth: request.DateOfBirth,
	}
}

func personResponseFromService(person service.Person) PersonResponse {
	return PersonResponse{
		ExternalID:  person.ExternalID,
		Name:        person.Name,
		Email:       person.Email,
		DateOfBirth: person.DateOfBirth,
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, ErrorResponse{
		Error:   code,
		Message: message,
	})
}

func writeValidationError(w http.ResponseWriter, message string, fields map[string]string) {
	writeJSON(w, http.StatusBadRequest, ValidationErrorResponse{
		Error:   errorValidation,
		Message: message,
		Fields:  fields,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
