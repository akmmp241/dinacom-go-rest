package app

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/controllers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"github.com/gofiber/fiber/v2"
	"time"
)

func NewRouter(authController controllers.AuthController) *fiber.App {
	appRouter := fiber.New(fiber.Config{
		Prefork:      true,
		AppName:      "Evia-BE-REST",
		IdleTimeout:  10 * time.Minute,
		ErrorHandler: exceptions.HandleError,
	})

	api := appRouter.Group("/api")

	auth := api.Group("/auth")

	auth.Post("/register", authController.Register)

	return appRouter
}
