package models

import (
	"github.com/google/uuid"
)

type Reset struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string
	Token     string `gorm:"unique"`
	ExpiresAt int64  // Unix timestamp in milliseconds
	Used      bool   `gorm:"default:false"`
}
