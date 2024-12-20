package repository

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"context"
	"database/sql"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (u UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	query := `INSERT INTO users (id, name, email, password) VALUES (NULL, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	user.Id = int(id)
	return user, nil
}

func (u UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = ?`
	rows, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}
	defer rows.Close()

	var user model.User
	if !rows.Next() {
		return nil, exceptions.NewNotFoundError()
	}

	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	return &user, nil
}
