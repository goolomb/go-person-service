package httpapi

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
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

	response := PersonResponse{
		ExternalID:  request.ExternalID,
		Name:        request.Name,
		Email:       request.Email,
		DateOfBirth: request.DateOfBirth,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
