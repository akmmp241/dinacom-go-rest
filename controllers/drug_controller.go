package controllers

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/gofiber/fiber/v2"
)

type DrugController interface {
	GetById(ctx *fiber.Ctx) error
}

type DrugControllerImpl struct {
	service.DrugService
}

func NewDrugController(drugService service.DrugService) *DrugControllerImpl {
	return &DrugControllerImpl{DrugService: drugService}
}

func (d DrugControllerImpl) GetById(ctx *fiber.Ctx) error {
	drugId := ctx.Params("drugId")
	if drugId == "" {
		return exceptions.NewBadRequestError("Drug id is required")
	}

	resp, err := d.DrugService.GetById(ctx.Context(), drugId)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Get drug success",
		Data:    resp,
		Errors:  nil,
	}

	return ctx.JSON(&globalResponse)
}
