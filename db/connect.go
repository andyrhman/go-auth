package db

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"go-auth/models"
)

var DB *gorm.DB

func Connect(){

	err := godotenv.Load()
	
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	// Connect to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = db

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&models.User{}, &models.Token{})

	log.Println("Connected to the database successfully!")
}