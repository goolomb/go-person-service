package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/goolomb/go-person-service/internal/service"
)

type PersonSaver interface {
	SavePerson(input service.SavePersonInput) (service.Person, map[string]string)
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

	person, validationErrors := h.personService.SavePerson(service.SavePersonInput{
		ExternalID:  request.ExternalID,
		Name:        request.Name,
		Email:       request.Email,
		DateOfBirth: request.DateOfBirth,
	})
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ValidationErrorResponse{
			Error:   "validation_error",
			Message: "Request body contains validation errors.",
			Fields:  validationErrors,
		})
		return
	}

	response := PersonResponse{
		ExternalID:  person.ExternalID,
		Name:        person.Name,
		Email:       person.Email,
		DateOfBirth: person.DateOfBirth,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
