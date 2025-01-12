package exceptions

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
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

func handleValidationErrorMessage(tag string, param string, field string) string {
	var msg string
	switch tag {
	case "required":
		msg = fmt.Sprintf("The %s field is required", field)
	case "email":
		msg = "This is not a valid email"
	case "min":
		msg = fmt.Sprintf("The %s field must be at least %s characters", strings.ToLower(param), tag)
	case "max":
		msg = fmt.Sprintf("The %s field must be at most %s characters", strings.ToLower(param), tag)
	case "eqfield":
		if param == "Password" {
			msg = "The password confirmation does not match"
		} else {
			msg = "The field does not match"
		}
	}

	return msg
}
