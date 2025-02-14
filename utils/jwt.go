package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID.String(), // Store UUID as string
		"exp": time.Now().Add(30 * time.Second).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_ACCESS")))
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID.String(), // Store UUID as string
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_REFRESH")))
}