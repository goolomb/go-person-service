package service

import (
	"errors"
	"testing"
	"time"

	"github.com/goolomb/go-person-service/internal/storage"
)

func TestSavePersonValidInputSucceedsAndTrimsFields(t *testing.T) {
	repo := &fakePersonRepository{}
	service := NewPersonService(repo)

	person, fields, err := service.SavePerson(SavePersonInput{
		ExternalID:  "3f93df6d-ff51-4740-9d27-fc6b2f30281c",
		Name:        " Jane Doe ",
		Email:       " jane@example.com ",
		DateOfBirth: "1990-01-02T03:04:05Z",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(fields) != 0 {
		t.Fatalf("expected no validation fields, got %v", fields)
	}
	if person.Name != "Jane Doe" {
		t.Fatalf("expected trimmed name, got %q", person.Name)
	}
	if person.Email != "jane@example.com" {
		t.Fatalf("expected trimmed email, got %q", person.Email)
	}
	if repo.created.Name != "Jane Doe" {
		t.Fatalf("expected repository name to be trimmed, got %q", repo.created.Name)
	}
	if repo.created.Email != "jane@example.com" {
		t.Fatalf("expected repository email to be trimmed, got %q", repo.created.Email)
	}
}

func TestSavePersonMissingRequiredFieldsReturnsAllValidationErrors(t *testing.T) {
	service := NewPersonService(&fakePersonRepository{})

	_, fields, err := service.SavePerson(SavePersonInput{})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}

	for _, field := range []string{"external_id", "name", "email", "date_of_birth"} {
		if fields[field] == "" {
			t.Fatalf("expected validation error for %q, got %v", field, fields)
		}
	}
}

func TestSavePersonInvalidFieldsReturnValidationErrors(t *testing.T) {
	service := NewPersonService(&fakePersonRepository{})

	_, fields, err := service.SavePerson(SavePersonInput{
		ExternalID:  "not-a-uuid",
		Name:        "Jane Doe",
		Email:       "not-an-email",
		DateOfBirth: "not-a-date",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}

	for _, field := range []string{"external_id", "email", "date_of_birth"} {
		if fields[field] == "" {
			t.Fatalf("expected validation error for %q, got %v", field, fields)
		}
	}
}

func TestSavePersonDuplicateRepositoryErrorMapsToServiceError(t *testing.T) {
	service := NewPersonService(&fakePersonRepository{createErr: storage.ErrAlreadyExists})

	_, _, err := service.SavePerson(validSavePersonInput())
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestGetPersonByExternalIDInvalidUUIDReturnsValidationError(t *testing.T) {
	service := NewPersonService(&fakePersonRepository{})

	_, fields, err := service.GetPersonByExternalID("not-a-uuid")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if fields["id"] == "" {
		t.Fatalf("expected validation error for id, got %v", fields)
	}
}

func TestGetPersonByExternalIDRepositoryNotFoundMapsToServiceError(t *testing.T) {
	service := NewPersonService(&fakePersonRepository{findErr: storage.ErrNotFound})

	_, _, err := service.GetPersonByExternalID("3f93df6d-ff51-4740-9d27-fc6b2f30281c")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type fakePersonRepository struct {
	created   storage.PersonModel
	createErr error
	findErr   error
}

func (r *fakePersonRepository) CreatePerson(person storage.PersonModel) (storage.PersonModel, error) {
	r.created = person
	if r.createErr != nil {
		return storage.PersonModel{}, r.createErr
	}

	person.ID = 1
	person.CreatedAt = time.Now().UTC()
	person.UpdatedAt = person.CreatedAt
	return person, nil
}

func (r *fakePersonRepository) FindPersonByExternalID(externalID string) (storage.PersonModel, error) {
	if r.findErr != nil {
		return storage.PersonModel{}, r.findErr
	}

	return storage.PersonModel{
		ID:          1,
		ExternalID:  externalID,
		Name:        "Jane Doe",
		Email:       "jane@example.com",
		DateOfBirth: time.Date(1990, time.January, 2, 3, 4, 5, 0, time.UTC),
	}, nil
}

func validSavePersonInput() SavePersonInput {
	return SavePersonInput{
		ExternalID:  "3f93df6d-ff51-4740-9d27-fc6b2f30281c",
		Name:        "Jane Doe",
		Email:       "jane@example.com",
		DateOfBirth: "1990-01-02T03:04:05Z",
	}
}
