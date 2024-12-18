package app

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:     true,
		AppName:     "Evia-BE-REST",
		IdleTimeout: 10 * time.Minute,
	})

	return app
}
