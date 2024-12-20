package controllers

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/gofiber/fiber/v2"
)

type AIController interface {
	Simplifier(ctx *fiber.Ctx) error
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

func NewAIController(AIService service.AIService) *AIControllerImpl {
	return &AIControllerImpl{AIService: AIService}
}
