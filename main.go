package main

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/app"
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/controllers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/middleware"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/go-playground/validator/v10"
)

func main() {
	cnf := config.NewConfig()

	db := app.NewDB(cnf)

	validate := validator.New()

	aiClient := config.InitAiClient(cnf)

	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()
	complaintRepo := repository.NewComplaintRepository()

	authService := service.NewAuthService(userRepo, sessionRepo, db, validate, cnf)
	aiService := service.NewComplaintService(validate, cnf, aiClient, complaintRepo, db)

	authController := controllers.NewAuthController(authService)
	aiController := controllers.NewComplaintController(aiService)

	mw := middleware.NewMiddleware(cnf, sessionRepo, userRepo, db)

	fiberApp := app.NewRouter(mw, authController, aiController)

	if err := fiberApp.Listen(":3000"); err != nil {
		panic(err)
	}
}
