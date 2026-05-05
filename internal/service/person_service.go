package service

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SavePersonInput struct {
	ExternalID  string
	Name        string
	Email       string
	DateOfBirth string
}

type Person struct {
	ExternalID  string
	Name        string
	Email       string
	DateOfBirth string
}

type PersonService struct{}

func (s PersonService) SavePerson(input SavePersonInput) (Person, map[string]string) {
	fields := make(map[string]string)

	externalID := strings.TrimSpace(input.ExternalID)
	if externalID == "" {
		fields["external_id"] = "external_id is required"
	} else if _, err := uuid.Parse(externalID); err != nil {
		fields["external_id"] = "external_id must be a valid UUID"
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		fields["name"] = "name is required"
	}

	email := strings.TrimSpace(input.Email)
	if email == "" {
		fields["email"] = "email is required"
	} else if !isValidEmail(email) {
		fields["email"] = "email must be a valid email address"
	}

	dateOfBirth := strings.TrimSpace(input.DateOfBirth)
	if dateOfBirth == "" {
		fields["date_of_birth"] = "date_of_birth is required"
	} else if _, err := time.Parse(time.RFC3339, dateOfBirth); err != nil {
		fields["date_of_birth"] = "date_of_birth must be a valid RFC3339 timestamp"
	}

	if len(fields) > 0 {
		return Person{}, fields
	}

	return Person{
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

	address, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if address.Name != "" || address.Address != email {
		return false
	}

	_, domain, found := strings.Cut(address.Address, "@")
	return found && strings.Contains(domain, ".")
}
