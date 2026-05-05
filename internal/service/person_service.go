package service

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/goolomb/go-person-service/internal/storage"
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

var (
	ErrValidation    = errors.New("validation error")
	ErrAlreadyExists = errors.New("person already exists")
	ErrNotFound      = errors.New("person not found")
)

type PersonRepository interface {
	CreatePerson(person storage.PersonModel) (storage.PersonModel, error)
	FindPersonByExternalID(externalID string) (storage.PersonModel, error)
}

type PersonService struct {
	repository PersonRepository
}

func NewPersonService(repo PersonRepository) *PersonService {
	return &PersonService{repository: repo}
}

func (s PersonService) SavePerson(input SavePersonInput) (Person, map[string]string, error) {
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
	var parsedDateOfBirth time.Time
	if dateOfBirth == "" {
		fields["date_of_birth"] = "date_of_birth is required"
	} else if parsed, err := time.Parse(time.RFC3339, dateOfBirth); err != nil {
		fields["date_of_birth"] = "date_of_birth must be a valid RFC3339 timestamp"
	} else {
		parsedDateOfBirth = parsed
	}

	if len(fields) > 0 {
		return Person{}, fields, ErrValidation
	}

	created, err := s.repository.CreatePerson(storage.PersonModel{
		ExternalID:  externalID,
		Name:        name,
		Email:       email,
		DateOfBirth: parsedDateOfBirth,
	})
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return Person{}, nil, ErrAlreadyExists
		}

		return Person{}, nil, err
	}

	return personFromModel(created), nil, nil
}

func (s PersonService) GetPersonByExternalID(id string) (Person, map[string]string, error) {
	externalID := strings.TrimSpace(id)
	fields := make(map[string]string)
	if externalID == "" {
		fields["id"] = "id is required"
	} else if _, err := uuid.Parse(externalID); err != nil {
		fields["id"] = "id must be a valid UUID"
	}

	if len(fields) > 0 {
		return Person{}, fields, ErrValidation
	}

	person, err := s.repository.FindPersonByExternalID(externalID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return Person{}, nil, ErrNotFound
		}

		return Person{}, nil, err
	}

	return personFromModel(person), nil, nil
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

func personFromModel(person storage.PersonModel) Person {
	return Person{
		ExternalID:  person.ExternalID,
		Name:        person.Name,
		Email:       person.Email,
		DateOfBirth: person.DateOfBirth.Format(time.RFC3339),
	}
}
