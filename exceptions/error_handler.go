package exceptions

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

var globalResponse = &model.GlobalResponse{
	Data: nil,
}

func HandleError(c *fiber.Ctx, err error) error {
	if errors.As(err, &validator.ValidationErrors{}) {
		return validationError(c, err.(validator.ValidationErrors))
	}

	if errors.As(err, &HttpInternalServerError{}) {
		log.Println(err.Error())
		return internalServerError(c)
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		globalResponse.Message = err.Error()
		return c.Status(e.Code).JSON(&globalResponse)
	}

	return baseError(c, err.(GlobalError))
}

func validationError(c *fiber.Ctx, err validator.ValidationErrors) error {
	globalResponse.Message = "Failed Validation"
	globalResponse.Errors = err.Error()

	log.Println(err.Error())
	return c.Status(fiber.StatusBadRequest).JSON(&globalResponse)
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
