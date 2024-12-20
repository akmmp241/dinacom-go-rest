package controllers

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type AuthControllerImpl struct {
	AuthService service.AuthService
}

func (con *AuthControllerImpl) Register(c *fiber.Ctx) error {
	registerRequest := &model.RegisterRequest{}
	err := c.BodyParser(registerRequest)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	registerResponse, err := con.AuthService.Register(c.Context(), *registerRequest)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Register success",
		Data:    &registerResponse,
		Errors:  nil,
	}

	return c.JSON(&globalResponse)
}

func (con *AuthControllerImpl) Login(c *fiber.Ctx) error {
	loginRequest := &model.LoginRequest{}
	err := c.BodyParser(loginRequest)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	loginResponse, err := con.AuthService.Login(c.Context(), *loginRequest)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Login success",
		Data:    &loginResponse,
		Errors:  nil,
	}

	return c.JSON(&globalResponse)
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{AuthService: authService}
}
