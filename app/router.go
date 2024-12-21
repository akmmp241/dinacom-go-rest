package app

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/controllers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/middleware"
	"github.com/gofiber/fiber/v2"
	"time"
)

func NewRouter(
	middleware middleware.Middleware,
	authController controllers.AuthController,
	aiController controllers.AIController,
) *fiber.App {
	appRouter := fiber.New(fiber.Config{
		Prefork:      true,
		AppName:      "Evia-BE-REST",
		IdleTimeout:  10 * time.Minute,
		ErrorHandler: exceptions.HandleError,
	})

	api := appRouter.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/me", authController.Me)

	ai := api.Use(middleware.Authenticate).Group("/ai")
	ai.Post("/simplify", aiController.Simplifier)
	ai.Post("/external/wound", aiController.ExternalWound)

	return appRouter
}
