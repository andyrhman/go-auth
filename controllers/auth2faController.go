package controllers

import (
	"encoding/base64"
	"encoding/base32"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
	// "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
)

type TwoFactorRequest struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Secret     string `json:"secret"`
	RememberMe bool   `json:"rememberMe"`
}

func TwoFactor(c *fiber.Ctx) error {
	var req TwoFactorRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Validate UUID
	if _, err := uuid.Parse(req.ID); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Find user
	var user models.User
	if err := db.DB.Where("id = ?", req.ID).First(&user).Error; err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Get secret
	secret := user.TFASecret
	if secret == "" {
		secret = req.Secret
	}

	// Verify code
	valid := totp.Validate(req.Code, secret)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Save secret if new
	if user.TFASecret == "" {
		if err := db.DB.Model(&user).Update("tfa_secret", secret).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving secret"})
		}
	}

	// Generate tokens
	userID, _ := uuid.Parse(req.ID)
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
	}

	// Set expiration based on rememberMe
	var expiration time.Time
	if req.RememberMe {
		expiration = time.Now().Add(365 * 24 * time.Hour) // 1 year
	} else {
		expiration = time.Now().Add(30 * time.Second) // 30 seconds
	}

	// Save refresh token
	refreshTokenRecord := models.Token{
		User_id:   userID,
		Token:     refreshToken,
		ExpiredAt: expiration,
	}

	if err := db.DB.Create(&refreshTokenRecord).Error; err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving token"})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  expiration,
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"token": accessToken})
}

// !! Fix this issue, it return the same secret key on 2fas auth app
func QR(c *fiber.Ctx) error {
	// Decode the base32 secret correctly 
	secretStr := "YRBGFHTE7J53MIVWE64S4HFRU2IPZ5PY"
	decodedSecret, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secretStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error decoding secret")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Go Auth",
		AccountName: "test2@mail.com",
		Secret:      decodedSecret, // Use decoded bytes
		SecretSize:  uint(len(decodedSecret) * 8), // Calculate proper size
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating QR code")
	}

	// Generate QR code
	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating QR code")
	}

	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(png)
	return c.SendString(fmt.Sprintf("<img src=\"data:image/png;base64,%s\" />", encoded))
}