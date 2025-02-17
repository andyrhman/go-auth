package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-auth/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.AuthenticatedUser)
	app.Post("/api/refresh", controllers.Refresh)
	app.Post("/api/logout", controllers.Logout)
	app.Post("/api/forgot", controllers.ForgotPassword)
	app.Post("/api/reset", controllers.ResetPassword)
	app.Post("/api/two-factor", controllers.TwoFactor)
	app.Get("/api/test", controllers.QR)
}
