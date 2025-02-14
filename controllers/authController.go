package controllers

import (
	"go-auth/db"
	"go-auth/models"

	"github.com/gofiber/fiber/v2"
	"go-auth/utils"
	"time"
	// "github.com/golang-jwt/jwt/v5"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	// Parse JSON body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	// Validate password confirmation
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// Hash password using utils.HashPassword
	hashedPassword := utils.HashPassword(data["password"])

	// Save user to database
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  []byte(hashedPassword),
	}

	db.DB.Create(user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parse JSON body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User

	if err := db.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Verify password
	if !utils.VerifyPassword(string(user.Password), data["password"]) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
	}

	// Save refresh token
	expiredAt := time.Now().Add(7 * 24 * time.Hour)
	refreshTokenRecord := models.Token{
		User_id:   user.Id,
		Token:     refreshToken,
		ExpiredAt: expiredAt,
	}

	if err := db.DB.Create(&refreshTokenRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving token"})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  expiredAt,
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"token": accessToken})
}
