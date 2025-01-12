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
	auth.Post("/google/callback", authController.GoogleCallback)
	auth.Get("/me", authController.Me)
	auth.Post("/forget/password", middleware.SendOtpMailRateLimiter, authController.ForgetPassword)
	auth.Post("/forget/password/verify", authController.VerifyForgetPasswordOtp)
	auth.Post("/reset/password", authController.ResetPassword)

	complaint := api.Group("/complaints")
	complaint.Use(middleware.Authenticate)
	complaint.Post("/", complaintController.ExternalWound)
	complaint.Get("/", complaintController.GetAll)
	complaint.Post("/simplify", complaintController.Simplifier)
	complaint.Get("/:complaintId", complaintController.GetById)
	complaint.Put("/:complaintId", complaintController.Update)
	complaint.Get("/:complaintId/recommendations", complaintController.GetRecommendedDrugs)

	return appRouter
}
