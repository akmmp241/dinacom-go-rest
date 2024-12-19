package repository

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"context"
	"database/sql"
)

type SessionRepository interface {
	Save(ctx context.Context, tx *sql.Tx, session *model.Session) (*model.Session, error)
}

type SessionRepositoryImpl struct {
}

func NewSessionRepository() *SessionRepositoryImpl {
	return &SessionRepositoryImpl{}
}

func (s SessionRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, session *model.Session) (*model.Session, error) {
	query := `INSERT INTO sessions (id, user_id, token, expires_at) VALUES (NULL, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, query, session.UserId, session.Token, session.ExpiresAt)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	session.Id = int(id)
	return session, nil
}
