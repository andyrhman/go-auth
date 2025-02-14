package utils

import (
	"encoding/base64"
	"fmt"
	"crypto/rand"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Generate a random salt (16 bytes)
func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

// HashPassword hashes a password using Argon2id and returns it in Argon2 standard format
func HashPassword(password string) string {
	salt := generateSalt()

	// Hash password using Argon2id
	hashedPassword := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	// Encode salt & hash in base64
	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hashedPassword)

	// Return password in Argon2 standard format
	return fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s", encodedSalt, encodedHash)
}

// VerifyPassword checks if the input password matches the stored Argon2 hash
func VerifyPassword(storedHash, inputPassword string) bool {
	parts := strings.Split(storedHash, "$")
	if len(parts) != 6 {
		return false // Invalid hash format
	}

	// Extract salt & stored hash
	encodedSalt := parts[4]
	encodedStoredHash := parts[5]

	// Decode salt & hash from base64
	salt, err1 := base64.StdEncoding.DecodeString(encodedSalt)
	_, err2 := base64.StdEncoding.DecodeString(encodedStoredHash)
	if err1 != nil || err2 != nil {
		return false
	}

	// Hash the input password using the extracted salt
	newHash := argon2.IDKey([]byte(inputPassword), salt, 3, 64*1024, 4, 32)

	// Compare the hashes
	return base64.StdEncoding.EncodeToString(newHash) == encodedStoredHash
}
