package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	User_id   uuid.UUID `gorm:"type:uuid"` // Changed to UUID type
	Token     string
	ExpiredAt time.Time
}
