package service

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"context"
	"database/sql"
	"errors"
)

type DrugService interface {
	GetById(ctx context.Context, id string) (*model.GetDrugDetailResponse, error)
}

type DrugServiceImpl struct {
	DB       *sql.DB
	DrugRepo repository.DrugRepository
}

func NewDrugService(drugRepo repository.DrugRepository, DB *sql.DB) *DrugServiceImpl {
	return &DrugServiceImpl{DrugRepo: drugRepo, DB: DB}
}

func (d DrugServiceImpl) GetById(ctx context.Context, id string) (*model.GetDrugDetailResponse, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	drug, err := d.DrugRepo.FindById(ctx, tx, id)
	if err != nil && errors.Is(err, exceptions.NotFoundError{}) {
		return nil, exceptions.NewHttpNotFoundError("Drug not found")
	} else if err != nil && !errors.Is(err, exceptions.NotFoundError{}) {
		return nil, err
	}

	_ = tx.Commit()

	getDrugDetailResponse := &model.GetDrugDetailResponse{
		Id:          drug.Id,
		BrandName:   drug.BrandName,
		Name:        drug.Name,
		Description: drug.Description,
		Price:       drug.Price,
		ImageUrl:    drug.ImageUrl,
	}

	return getDrugDetailResponse, nil
}
