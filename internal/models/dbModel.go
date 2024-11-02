package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Role     string `gorm:"not null"`
	Contact  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
