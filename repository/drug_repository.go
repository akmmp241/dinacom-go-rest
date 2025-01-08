package repository

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"context"
	"database/sql"
)

type DrugRepository interface {
	Save(ctx context.Context, tx *sql.Tx, drug *model.Drug) (*model.Drug, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]model.Drug, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Drug, error)
	Search(ctx context.Context, tx *sql.Tx, name string) ([]model.Drug, error)
}

type DrugRepositoryImpl struct {
}

func NewDrugRepository() *DrugRepositoryImpl {
	return &DrugRepositoryImpl{}
}

func (d DrugRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, drug *model.Drug) (*model.Drug, error) {
	query := "INSERT INTO drugs (id, brand_name, name, description, price, image_url) VALUES (NULL, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, drug.BrandName, drug.Name, &drug.Price, &drug.Description, drug.ImageUrl)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	drug.Id = int(id)
	return drug, nil
}

func (d DrugRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]model.Drug, error) {
	query := "SELECT id, brand_name, name, price, description, image_url FROM drugs"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}
	defer rows.Close()

	var drugs []model.Drug
	for rows.Next() {
		var drug model.Drug
		err := rows.Scan(&drug.Id, &drug.BrandName, &drug.Name, &drug.Price, &drug.Description, &drug.ImageUrl)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		drugs = append(drugs, drug)
	}

	return drugs, nil
}

func (d DrugRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Drug, error) {
	query := "SELECT id, brand_name, name, price, description, image_url FROM drugs WHERE id = ?"
	row := tx.QueryRowContext(ctx, query, id)

	var drug model.Drug
	err := row.Scan(&drug.Id, &drug.BrandName, &drug.Name, &drug.Price, &drug.Description, &drug.ImageUrl)
	if err != nil {
		return nil, exceptions.NewNotFoundError()
	}

	return &drug, nil
}

func (d DrugRepositoryImpl) Search(ctx context.Context, tx *sql.Tx, name string) ([]model.Drug, error) {
	query := "SELECT id, brand_name, name, price, description, image_url FROM drugs WHERE name LIKE ?"
	rows, err := tx.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}
	defer rows.Close()

	var drugs []model.Drug
	for rows.Next() {
		var drug model.Drug
		err := rows.Scan(&drug.Id, &drug.BrandName, &drug.Name, &drug.Price, &drug.Description, &drug.ImageUrl)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		drugs = append(drugs, drug)
	}

	return drugs, nil
}
