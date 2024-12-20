package main

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/app"
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/controllers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/go-playground/validator/v10"
)

func main() {
	cnf := config.NewConfig()

	db := app.NewDB(cnf)

	validate := validator.New()

	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()

	authService := service.NewAuthService(userRepo, sessionRepo, db, validate, cnf)

	authController := controllers.NewAuthController(authService)

	fiberApp := app.NewRouter(authController)

	if err := fiberApp.Listen(":3000"); err != nil {
		panic(err)
	}
}
