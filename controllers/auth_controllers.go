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
	Me(c *fiber.Ctx) error
	ForgetPassword(c *fiber.Ctx) error
	VerifyForgetPasswordOtp(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
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

func (con *AuthControllerImpl) Me(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return exceptions.NewBadRequestError("Missing access token")
	}

	meResponse, err := con.AuthService.Me(c.Context(), token)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Get user success",
		Data:    meResponse,
		Errors:  nil,
	}

	return c.JSON(&globalResponse)
}

func (con *AuthControllerImpl) ForgetPassword(c *fiber.Ctx) error {
	forgetPasswordRequest := &model.ForgetPasswordRequest{}
	err := c.BodyParser(forgetPasswordRequest)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	err = con.AuthService.ForgetPassword(c.Context(), *forgetPasswordRequest)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Send otp email success",
		Data:    nil,
		Errors:  nil,
	}

	return c.JSON(&globalResponse)
}

func (con *AuthControllerImpl) VerifyForgetPasswordOtp(c *fiber.Ctx) error {
	req := &model.VerifyForgetPasswordOtpRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	verifyForgetPasswordOtpResponse, err := con.AuthService.VerifyForgetPasswordOtp(c.Context(), *req)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Forget password otp verified",
		Data:    verifyForgetPasswordOtpResponse,
		Errors:  nil,
	}

	return c.JSON(&globalResponse)
}

func (con *AuthControllerImpl) ResetPassword(c *fiber.Ctx) error {
	resetPasswordRequest := &model.ResetPasswordRequest{}
	err := c.BodyParser(resetPasswordRequest)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	message, err := con.AuthService.ResetPassword(c.Context(), *resetPasswordRequest)
	if err != nil {
		return err
	}
	globalResponse := model.GlobalResponse{
		Message: message.Message,
		Data:    nil,
		Errors:  nil,
	}

	return c.JSON(&globalResponse)
}

func (con *AuthControllerImpl) GoogleCallback(c *fiber.Ctx) error {
	req := &model.GoogleCallbackRequest{}
	err := c.BodyParser(req)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	loginResponse, err := con.AuthService.GoogleCallback(c.Context(), *req)
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
