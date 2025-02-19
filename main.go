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
	redis := app.NewRedisClient(cnf)
	mailer := config.NewMailer(cnf)

	validate := validator.New()

	aiClient := config.InitAiClient(cnf)
	awsClient := config.InitS3Client(cnf)
	oauthClient := config.NewOauthClient(cnf)

	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()
	complaintRepo := repository.NewComplaintRepository()
	drugRepo := repository.NewDrugRepository()

	authService := service.NewAuthService(userRepo, sessionRepo, db, validate, cnf, redis, mailer, oauthClient)
	complaintService := service.NewComplaintService(validate, cnf, aiClient, awsClient, complaintRepo, db, drugRepo)
	drugService := service.NewDrugService(drugRepo, db)

	authController := controllers.NewAuthController(authService)
	complaintController := controllers.NewComplaintController(complaintService)
	drugController := controllers.NewDrugController(drugService)

	mw := middleware.NewMiddleware(cnf, sessionRepo, userRepo, db, redis)

	fiberApp := app.NewRouter(mw, authController, complaintController, drugController)

	if err := fiberApp.Listen(":3000"); err != nil {
		panic(err)
	}
}
