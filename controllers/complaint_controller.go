package controllers

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/service"
	"github.com/gofiber/fiber/v2"
)

type ComplaintController interface {
	Simplifier(ctx *fiber.Ctx) error
	ExternalWound(ctx *fiber.Ctx) error
	GetById(ctx *fiber.Ctx) error
	GetAll(ctx *fiber.Ctx) error
}

type ComplaintControllerImpl struct {
	ComplaintService service.ComplaintService
}

func (A ComplaintControllerImpl) Simplifier(ctx *fiber.Ctx) error {
	simplifyRequest := &model.SimplifyRequest{}
	err := ctx.BodyParser(simplifyRequest)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	resp, err := A.ComplaintService.Simplifier(ctx.Context(), simplifyRequest)
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

func (A ComplaintControllerImpl) ExternalWound(ctx *fiber.Ctx) error {
	req := &model.ComplaintRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return exceptions.NewBadRequestError("Invalid request body")
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return exceptions.NewBadRequestError("Image is required")
	}

	open, err := file.Open()
	if err != nil {
		return exceptions.NewInternalServerError()
	}
	defer open.Close()
	req.Image = open

	user := ctx.UserContext().Value("user").(*model.User)

	resp, err := A.ComplaintService.ExternalWound(ctx.Context(), req, user)
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

func (A ComplaintControllerImpl) GetById(ctx *fiber.Ctx) error {
	complaintId := ctx.Params("complaintId")
	if complaintId == "" {
		return exceptions.NewBadRequestError("Complaint id is required")
	}

	user := ctx.UserContext().Value("user").(*model.User)

	resp, err := A.ComplaintService.GetById(ctx.Context(), complaintId, user)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Get complaint success",
		Data:    resp,
		Errors:  nil,
	}

	return ctx.JSON(&globalResponse)
}

func (A ComplaintControllerImpl) GetAll(ctx *fiber.Ctx) error {
	user := ctx.UserContext().Value("user").(*model.User)

	resp, err := A.ComplaintService.GetAll(ctx.Context(), user)
	if err != nil {
		return err
	}

	globalResponse := model.GlobalResponse{
		Message: "Get all complaint success",
		Data:    resp,
		Errors:  nil,
	}

	return ctx.JSON(&globalResponse)
}

func NewComplaintController(ComplaintService service.ComplaintService) *ComplaintControllerImpl {
	return &ComplaintControllerImpl{ComplaintService: ComplaintService}
}
