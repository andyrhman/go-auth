package main

import (
	"github.com/gofiber/fiber/v2"
	"go-auth/database"
	"go-auth/routes"
)

func main() {
	database.Connect()

    app := fiber.New()

	routes.Setup(app)
	
    app.Listen(":8000")
}
