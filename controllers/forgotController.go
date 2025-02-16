package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"net/smtp"
)

func ForgotPassword(c *fiber.Ctx) error {
	type ForgotInput struct {
		Email string `json:"email" validate:"required,email"`
	}

	input := new(ForgotInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Generate random token
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
	}
	tokenStr := hex.EncodeToString(token)

	// Save reset token
	resetRecord := models.Reset{
		Email:     input.Email,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(30 * time.Minute).UnixMilli(),
	}

	if err := db.DB.Create(&resetRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving reset token"})
	}

	// Send email
	if err := sendResetEmail(input.Email, tokenStr); err != nil {
		fmt.Println("Failed to send email:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error sending email"})
	}

	return c.JSON(fiber.Map{"message": "Please check your email"})
}

func ResetPassword(c *fiber.Ctx) error {
	type ResetInput struct {
		Token           string `json:"token"`
		Password        string `json:"password" validate:"required,min=6"`
		PasswordConfirm string `json:"password_confirm" validate:"required,min=6"`
	}

	input := new(ResetInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if input.Password != input.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Passwords do not match"})
	}

	// Find reset token
	var resetToken models.Reset
	if err := db.DB.Where("token = ?", input.Token).First(&resetToken).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid token"})
	}

	if resetToken.Used || resetToken.ExpiresAt < time.Now().UnixMilli() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Token expired or already used"})
	}

	// Find user
	var user models.User
	if err := db.DB.Where("email = ?", resetToken.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Update password
	hashedPassword := utils.HashPassword(input.Password)
	if err := db.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating password"})
	}

	// Mark token as used
	if err := db.DB.Model(&resetToken).Update("used", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating token"})
	}

	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}

func sendResetEmail(email, token string) error {
	// Parse template
	url := fmt.Sprintf("http://%s/reset/%s", os.Getenv("APP_HOST"), token)
	html, err := utils.ParseTemplate("templates/forgot.html", struct {
		Email string
		URL   string
	}{
		Email: email,
		URL:   url,
	})
	if err != nil {
		return err
	}

	// SMTP configuration
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	// Create authentication
	// auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// MailHog
	var auth smtp.Auth
	if smtpUser != "" && smtpPass != "" {
		auth = smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	}
	
	// Email headers
	subject := "Subject: Reset Your Password\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	to := []string{email}
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + email + "\r\n" +
			subject +
			mime +
			html,
	)

	// Send email
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		to,
		msg,
	)

	return err
}
