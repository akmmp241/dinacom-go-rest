package exceptions

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

type IError struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

var globalResponse = &model.GlobalResponse{
	Data: nil,
}

func HandleError(c *fiber.Ctx, err error) error {
	if errors.As(err, &FailedValidationError{}) {
		return validationError(c, err.(FailedValidationError))
	}

	if errors.As(err, &HttpInternalServerError{}) {
		log.Println(err.Error())
		return internalServerError(c)
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		globalResponse.Message = e.Message
		globalResponse.Errors = nil
		return c.Status(e.Code).JSON(&globalResponse)
	}

	return baseError(c, err.(GlobalError))
}

func validationError(c *fiber.Ctx, err FailedValidationError) error {
	globalResponse.Message = err.Msg
	globalResponse.Errors = err.Errors

	log.Println(err.Error())
	return c.Status(err.GetCode()).JSON(&globalResponse)
}

func baseError(c *fiber.Ctx, err GlobalError) error {
	globalResponse.Message = err.Error()
	globalResponse.Errors = nil

	log.Println(err.Error())
	return c.Status(err.GetCode()).JSON(&globalResponse)
}

func internalServerError(c *fiber.Ctx) error {
	globalResponse.Message = "Internal Server Error"
	globalResponse.Errors = nil

	return c.Status(fiber.StatusInternalServerError).JSON(&globalResponse)
}
