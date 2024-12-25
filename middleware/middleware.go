package middleware

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"akmmp241/dinamcom-2024/dinacom-go-rest/exceptions"
	"akmmp241/dinamcom-2024/dinacom-go-rest/helpers"
	"akmmp241/dinamcom-2024/dinacom-go-rest/repository"
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"time"
)

type Middleware interface {
	Authenticate(c *fiber.Ctx) error
	SendOtpMailRateLimiter(c *fiber.Ctx) error
}

type MiddlewareImpl struct {
	SessionRepo repository.SessionRepository
	UserRepo    repository.UserRepository
	Cnf         *config.Config
	DB          *sql.DB
	RedisClient *redis.Client
}

func NewMiddleware(
	cnf *config.Config,
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
	db *sql.DB,
	redisClient *redis.Client,
) *MiddlewareImpl {
	return &MiddlewareImpl{
		Cnf:         cnf,
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		DB:          db,
		RedisClient: redisClient,
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

	user, err := i.UserRepo.FindById(c.Context(), tx, session.UserId)
	if err != nil {
		return err
	}

	_ = tx.Commit()

	ctx := context.WithValue(c.UserContext(), "user", user)
	c.SetUserContext(ctx)

	return c.Next()
}

func (i *MiddlewareImpl) SendOtpMailRateLimiter(c *fiber.Ctx) error {
	ip := c.IP()
	key := "send_otp_mail:" + ip
	err := i.RedisClient.Get(c.Context(), key).Err()
	if err == nil {
		return &fiber.Error{
			Code:    fiber.StatusTooManyRequests,
			Message: "Too many requests",
		}
	}

	err = i.RedisClient.Set(c.Context(), key, 1, 1*time.Minute).Err()
	if err != nil {
		return exceptions.NewInternalServerError()
	}

	return c.Next()
}
