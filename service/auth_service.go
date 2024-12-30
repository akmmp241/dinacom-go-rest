package service

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/helpers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/model"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"html/template"
	"log"
	"time"
)

//go:embed mail-templates/send-otp.html
var OTPTemplateEmail string

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	Me(ctx context.Context, token string) (*model.MeResponse, error)
	ForgetPassword(ctx context.Context, req model.ForgetPasswordRequest) error
	VerifyForgetPasswordOtp(ctx context.Context, req model.VerifyForgetPasswordOtpRequest) (*model.VerifyForgetPasswordOtpResponse, error)
	ResetPassword(ctx context.Context, req model.ResetPasswordRequest) (*model.ResetPasswordResponse, error)
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
		Email: user.Email,
	}, nil
}

func (s AuthServiceImpl) ForgetPassword(ctx context.Context, req model.ForgetPasswordRequest) error {
	err := s.Validate.Struct(req)
	if err != nil {
		return exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return exceptions.NewInternalServerError()
	}

	_, err = s.UserRepo.FindByEmail(ctx, tx, req.Email)
	if err != nil {
		time.Sleep(3 * time.Second)
		return nil
	}

	otpCode := helpers.GenerateRandomCodeForOtp()
	key := fmt.Sprintf("otp:%s", req.Email)
	err = s.RedisClient.SetEx(ctx, key, otpCode, time.Minute*5).Err()
	if err != nil {
		log.Println("error while set otp to redis", err)
		return exceptions.NewInternalServerError()
	}

	tmpl, err := template.New("email").Parse(OTPTemplateEmail)
	if err != nil {
		log.Println("error while parse template", err)
		return exceptions.NewInternalServerError()
	}

	otpData := config.SendOtpEmailData{
		OtpCode: otpCode,
		Email:   req.Email,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, otpData); err != nil {
		log.Println("error while execute template", err)
		return exceptions.NewInternalServerError()
	}

	err = s.Mailer.SendEmail(req.Email, "Forget Password OTP", body.String())
	if err != nil {
		log.Println("error while send email", err)
		return exceptions.NewInternalServerError()
	}

	return nil
}

func (s AuthServiceImpl) VerifyForgetPasswordOtp(ctx context.Context, req model.VerifyForgetPasswordOtpRequest) (*model.VerifyForgetPasswordOtpResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return nil, exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
	}

	key := fmt.Sprintf("otp:%s", req.Email)
	otpCode, err := s.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, exceptions.NewBadRequestError("Invalid OTP")
	}

	if otpCode != req.Otp {
		return nil, exceptions.NewBadRequestError("Invalid OTP")
	}

	err = s.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Println("error while delete otp from redis", err)
		return nil, exceptions.NewInternalServerError()
	}

	resetPasswordKey := fmt.Sprintf("reset-password:%s", req.Email)
	resetPasswordToken := uuid.NewString()
	encodedToken := helpers.StringToBase64([]byte(resetPasswordToken))
	encryptedToken, err := helpers.Encrypt(encodedToken, s.Cnf.Env.GetString("APP_KEY"))
	if err != nil {
		log.Println("error while encrypt token", err)
		return nil, exceptions.NewInternalServerError()
	}

	err = s.RedisClient.SetEx(ctx, resetPasswordKey, resetPasswordToken, time.Minute*5).Err()
	if err != nil {
		log.Println("error while set reset password token to redis", err)
		return nil, exceptions.NewInternalServerError()
	}

	verifyForgetPasswordOtpResponse := model.VerifyForgetPasswordOtpResponse{
		Email:              req.Email,
		ResetPasswordToken: encryptedToken,
	}

	return &verifyForgetPasswordOtpResponse, nil
}

func (s AuthServiceImpl) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) (*model.ResetPasswordResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return nil, exceptions.NewFailedValidationError(req, err.(validator.ValidationErrors))
	}

	hashPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}
	_, err = s.UserRepo.UpdatePassword(ctx, tx, req.Email, hashPassword)
	if err != nil {
		return nil, err
	}

	_ = tx.Commit()

	resetPasswordResponse := model.ResetPasswordResponse{
		Message: "Success Reset Password",
	}
	return &resetPasswordResponse, nil
}
