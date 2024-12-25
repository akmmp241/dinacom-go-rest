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
	complaintController controllers.ComplaintController,
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
	auth.Use(middleware.SendOtpMailRateLimiter).Post("/forget/password", authController.ForgetPassword)
	auth.Post("/forget/password/verify", authController.VerifyForgetPasswordOtp)
	auth.Post("/reset/password", authController.ResetPassword)

	authApi := api.Use(middleware.Authenticate)
	authApi.Post("/ai/simplify", complaintController.Simplifier)

	complaint := authApi.Group("/complaints")
	complaint.Post("/", complaintController.ExternalWound)
	complaint.Get("/", complaintController.GetAll)
	complaint.Get("/:complaintId", complaintController.GetById)

	return appRouter
}
