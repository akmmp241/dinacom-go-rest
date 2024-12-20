package service

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/helpers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error)
}

type AuthServiceImpl struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
	DB          *sql.DB
	Validate    *validator.Validate
	Cnf         *config.Config
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, DB *sql.DB, validate *validator.Validate, cnf *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{UserRepo: userRepo, SessionRepo: sessionRepo, DB: DB, Validate: validate, Cnf: cnf}
}

func (s AuthServiceImpl) Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error) {
	err := s.Validate.Struct(&req)
	if err != nil {
		return nil, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	user, err := s.UserRepo.FindByEmail(ctx, tx, req.Email)
	if err != nil && !errors.Is(err, exceptions.NotFoundError{}) {
		return nil, err
	}

	if user != nil {
		return nil, exceptions.NewBadRequestError("Email already registered")
	}

	hashedPassword, err := helpers.HashPassword(req.Password)

	user = &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	user, err = s.UserRepo.Save(ctx, tx, user)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	token := uuid.NewString()
	encodedToken := helpers.StringToBase64([]byte(token))
	encryptedToken, err := helpers.Encrypt(encodedToken, s.Cnf.Env.GetString("APP_KEY"))
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	_, err = s.SessionRepo.Save(ctx, tx, &model.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	})

	_ = tx.Commit()

	return &model.RegisterResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Token: encryptedToken,
	}, nil
}
