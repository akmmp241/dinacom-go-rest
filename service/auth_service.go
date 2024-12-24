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
	"github.com/redis/go-redis/v9"
	"html/template"
	"log"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	Me(ctx context.Context, token string) (*model.MeResponse, error)
}

type AuthServiceImpl struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
	DB          *sql.DB
	Validate    *validator.Validate
	Cnf         *config.Config
	RedisClient *redis.Client
	Mailer      *config.Mailer
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	DB *sql.DB, validate *validator.Validate,
	cnf *config.Config,
	redisClient *redis.Client,
	mailer *config.Mailer,
) *AuthServiceImpl {
	return &AuthServiceImpl{UserRepo: userRepo, SessionRepo: sessionRepo, DB: DB, Validate: validate, Cnf: cnf, RedisClient: redisClient, Mailer: mailer}
}

func (s AuthServiceImpl) Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error) {
	err := s.Validate.Struct(&req)
	if err != nil {
		return nil, exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
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
	if err != nil {
		_ = tx.Rollback()
		log.Println("here 2")
		return nil, exceptions.NewInternalServerError()
	}

	_ = tx.Commit()

	return &model.RegisterResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Token: encryptedToken,
	}, nil
}

func (s AuthServiceImpl) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	err := s.Validate.Struct(&req)
	if err != nil {
		return nil, exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	user, err := s.UserRepo.FindByEmail(ctx, tx, req.Email)
	if err != nil && errors.Is(err, exceptions.NotFoundError{}) {
		return nil, exceptions.NewHttpConflictError("Invalid Credentials")
	} else if err != nil && !errors.Is(err, exceptions.NotFoundError{}) {
		return nil, err
	}

	if !helpers.VerifyPassword(req.Password, user.Password) {
		return nil, exceptions.NewHttpConflictError("Invalid Credentials")
	}

	token := uuid.NewString()
	encodedToken := helpers.StringToBase64([]byte(token))
	encryptedToken, err := helpers.Encrypt(encodedToken, s.Cnf.Env.GetString("APP_KEY"))
	if err != nil {
		log.Println("here 1")
		return nil, exceptions.NewInternalServerError()
	}

	_, err = s.SessionRepo.Save(ctx, tx, &model.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	})
	if err != nil {
		_ = tx.Rollback()
		log.Println("here 2")
		return nil, exceptions.NewInternalServerError()
	}

	_ = tx.Commit()

	return &model.LoginResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Token: encryptedToken,
	}, nil
}

func (s AuthServiceImpl) Me(ctx context.Context, token string) (*model.MeResponse, error) {
	decodedToken, err := helpers.Decrypt(token, s.Cnf.Env.GetString("APP_KEY"))
	if err != nil {
		return nil, exceptions.NewUnauthorizedError("Unauthorized")
	}

	decodedTokenBase64 := helpers.Base64ToString(decodedToken)

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	session, err := s.SessionRepo.FindByToken(ctx, tx, decodedTokenBase64)
	if err != nil && errors.Is(err, exceptions.NotFoundError{}) {
		return nil, exceptions.NewUnauthorizedError("Unauthorized")
	} else if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, exceptions.NewUnauthorizedError("Unauthorized")
	}

	user, err := s.UserRepo.FindById(ctx, tx, session.UserId)
	if err != nil {
		return nil, err
	}

	_ = tx.Commit()

	return &model.MeResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
