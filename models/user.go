package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  []byte
}

