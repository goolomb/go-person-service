package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goolomb/go-person-service/internal/httpapi"
)

type PersonService struct{}

func (s PersonService) SavePerson(request httpapi.SavePersonRequest) (httpapi.PersonResponse, map[string]string) {
	fields := make(map[string]string)

	externalID := strings.TrimSpace(request.ExternalID)
	if externalID == "" {
		fields["external_id"] = "external_id is required"
	} else if _, err := uuid.Parse(externalID); err != nil {
		fields["external_id"] = "external_id must be a valid UUID"
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		fields["name"] = "name is required"
	}

	email := strings.TrimSpace(request.Email)
	if email == "" {
		fields["email"] = "email is required"
	} else if !isValidEmail(email) {
		fields["email"] = "email must be a valid email address"
	}

	dateOfBirth := strings.TrimSpace(request.DateOfBirth)
	if dateOfBirth == "" {
		fields["date_of_birth"] = "date_of_birth is required"
	} else if _, err := time.Parse(time.RFC3339, dateOfBirth); err != nil {
		fields["date_of_birth"] = "date_of_birth must be a valid RFC3339 timestamp"
	}

	if len(fields) > 0 {
		return httpapi.PersonResponse{}, fields
	}

	return httpapi.PersonResponse{
		ExternalID:  externalID,
		Name:        name,
		Email:       email,
		DateOfBirth: dateOfBirth,
	}, nil
}

func isValidEmail(email string) bool {
	if strings.ContainsAny(email, " \t\r\n") {
		return false
	}

	local, domain, found := strings.Cut(email, "@")
	if !found || strings.Contains(domain, "@") {
		return false
	}

	return local != "" && domain != "" && strings.Contains(domain, ".")
}
