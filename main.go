package main

import (
	"github.com/gofiber/fiber/v2"
	"go-auth/db"
	"go-auth/routes"
)

func main() {
	db.Connect()

    app := fiber.New()

	routes.Setup(app)
	
    app.Listen(":8000")
}
