package repository

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"context"
	"database/sql"
)

type ComplaintRepository interface {
	Save(ctx context.Context, tx *sql.Tx, complaints *model.Complaint) (*model.Complaint, error)
	FindAll(ctx context.Context, tx *sql.Tx, userId int) ([]model.Complaint, error)
	FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Complaint, error)
}

type ComplaintRepositoryImpl struct {
}

func NewComplaintRepository() *ComplaintRepositoryImpl {
	return &ComplaintRepositoryImpl{}
}

func (c ComplaintRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, complaints *model.Complaint) (*model.Complaint, error) {
	query := `INSERT INTO complaints (id, user_id, title, complaints, response, image_url, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, &complaints.Id, &complaints.UserId, &complaints.Title, &complaints.ComplaintsMsg, &complaints.Response, &complaints.ImageUrl, &complaints.CreatedAt)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	return complaints, nil
}

func (c ComplaintRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, userId int) ([]model.Complaint, error) {
	query := `SELECT id, user_id, title, complaints, response, image_url, created_at FROM complaints WHERE user_id = ?`
	rows, err := tx.QueryContext(ctx, query, &userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var complaints []model.Complaint
	for rows.Next() {
		var complaint model.Complaint
		err := rows.Scan(&complaint.Id, &complaint.UserId, &complaint.Title, &complaint.ComplaintsMsg, &complaint.Response, &complaint.ImageUrl, &complaint.CreatedAt)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

func (c ComplaintRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (*model.Complaint, error) {
	query := `SELECT id, user_id, title, complaints, response, image_url, created_at FROM complaints WHERE id = ?`
	rows, err := tx.QueryContext(ctx, query, &id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var complaint model.Complaint
	if !rows.Next() {
		return nil, exceptions.NewNotFoundError()
	}

	err = rows.Scan(&complaint.Id, &complaint.UserId, &complaint.Title, &complaint.ComplaintsMsg, &complaint.Response, &complaint.ImageUrl, &complaint.CreatedAt)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	return &complaint, nil
}
