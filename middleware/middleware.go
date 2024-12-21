package middleware

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/helpers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Middleware interface {
	Authenticate(c *fiber.Ctx) error
}

type MiddlewareImpl struct {
	SessionRepo repository.SessionRepository
	UserRepo    repository.UserRepository
	Cnf         *config.Config
	DB          *sql.DB
}

func NewMiddleware(
	cnf *config.Config,
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
	db *sql.DB,
) *MiddlewareImpl {
	return &MiddlewareImpl{
		Cnf:         cnf,
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		DB:          db,
	}
}

func (i *MiddlewareImpl) Authenticate(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		return exceptions.NewBadRequestError("Missing access token")
	}

	decodedToken, err := helpers.Decrypt(accessToken, i.Cnf.Env.GetString("APP_KEY"))
	if err != nil {
		return exceptions.NewUnauthorizedError("Unauthorized")
	}

	decodedTokenBase64 := helpers.Base64ToString(decodedToken)

	tx, err := i.DB.Begin()
	if err != nil {
		return exceptions.NewInternalServerError()
	}

	session, err := i.SessionRepo.FindByToken(c.Context(), tx, decodedTokenBase64)
	if err != nil && errors.Is(err, exceptions.NotFoundError{}) {
		return exceptions.NewUnauthorizedError("Unauthorized")
	} else if err != nil {
		return err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return exceptions.NewUnauthorizedError("Unauthorized")
	}

	_, err = i.UserRepo.FindById(c.Context(), tx, session.UserId)
	if err != nil {
		return err
	}

	_ = tx.Commit()

	return c.Next()
}
