package storage

import "time"

type PersonModel struct {
	ID          uint      `gorm:"primaryKey"`
	ExternalID  string    `gorm:"not null;uniqueIndex"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"not null"`
	DateOfBirth time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
