package storage

import (
	"errors"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("person not found")
	ErrAlreadyExists = errors.New("person already exists")
)

type PersonRepository struct {
	db *gorm.DB
}

func OpenSQLiteDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&PersonModel{}); err != nil {
		return nil, err
	}

	return db, nil
}

func NewPersonRepository(db *gorm.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) CreatePerson(person PersonModel) (PersonModel, error) {
	if err := r.db.Create(&person).Error; err != nil {
		if isUniqueConstraintError(err) {
			return PersonModel{}, ErrAlreadyExists
		}

		return PersonModel{}, err
	}

	return person, nil
}

func (r *PersonRepository) FindPersonByExternalID(externalID string) (PersonModel, error) {
	var person PersonModel
	if err := r.db.Where("external_id = ?", externalID).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PersonModel{}, ErrNotFound
		}

		return PersonModel{}, err
	}

	return person, nil
}

func isUniqueConstraintError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unique constraint failed") &&
		strings.Contains(message, "external_id")
}
