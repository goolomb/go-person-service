package storage

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestPersonRepository(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "people.db")
	db, err := OpenSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	repository := NewPersonRepository(db)
	dateOfBirth := time.Date(1990, time.January, 2, 3, 4, 5, 0, time.UTC)
	person := PersonModel{
		ExternalID:  "4ef4de59-8bc2-40dd-925d-6987a292006c",
		Name:        "Jane Doe",
		Email:       "jane@example.com",
		DateOfBirth: dateOfBirth,
	}

	created, err := repository.CreatePerson(person)
	if err != nil {
		t.Fatalf("create person: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("expected created person to have an ID")
	}

	found, err := repository.FindPersonByExternalID(person.ExternalID)
	if err != nil {
		t.Fatalf("find person by external ID: %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("expected ID %d, got %d", created.ID, found.ID)
	}
	if found.ExternalID != person.ExternalID {
		t.Fatalf("expected ExternalID %q, got %q", person.ExternalID, found.ExternalID)
	}
	if found.Name != person.Name {
		t.Fatalf("expected Name %q, got %q", person.Name, found.Name)
	}
	if found.Email != person.Email {
		t.Fatalf("expected Email %q, got %q", person.Email, found.Email)
	}
	if !found.DateOfBirth.Equal(dateOfBirth) {
		t.Fatalf("expected DateOfBirth %s, got %s", dateOfBirth, found.DateOfBirth)
	}

	_, err = repository.CreatePerson(person)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}

	_, err = repository.FindPersonByExternalID("d9bf4951-d35d-45d5-b6a9-31d0e879444d")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
