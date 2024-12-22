package controllers

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AIController interface {
	Simplifier(ctx *fiber.Ctx) error
	ExternalWound(ctx *fiber.Ctx) error
}

type AIControllerImpl struct {
	AIService service.AIService
}

func (A AIControllerImpl) Simplifier(ctx *fiber.Ctx) error {
	simplifyRequest := &model.SimplifyRequest{}
	err := ctx.BodyParser(simplifyRequest)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	resp, err := A.AIService.Simplifier(ctx.Context(), simplifyRequest)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Simplify success",
		Data:    resp,
		Errors:  nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(&globalResponse)
}

func (A AIControllerImpl) ExternalWound(ctx *fiber.Ctx) error {
	req := &model.ExternalWoundRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		log.Println("Failed to get file", err.Error())
		return exceptions.NewInternalServerError()
	}

	open, err := file.Open()
	if err != nil {
		return exceptions.NewInternalServerError()
	}
	defer open.Close()
	req.Image = open

	user := ctx.UserContext().Value("user").(*model.User)

	resp, err := A.AIService.ExternalWound(ctx.Context(), req, user)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Identify external wound success",
		Data:    resp,
		Errors:  nil,
	}

	return ctx.JSON(&globalResponse)
}

func NewAIController(AIService service.AIService) *AIControllerImpl {
	return &AIControllerImpl{AIService: AIService}
}
