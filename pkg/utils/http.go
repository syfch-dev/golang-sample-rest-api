package utils

import (
	"context"
	"my-go-project/internal/models"
	"my-go-project/pkg/httpErrors"

	"github.com/labstack/echo/v5"
)

func GetConfigName(env string) string {
	if env == "docker" {
		return "config-docker"
	}
	return "config-local"
}

type UserCtxKey struct{}

// Get request id from echo context
func GetRequestID(c *echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, httpErrors.Unauthorized
	}

	return user, nil
}

// Get user ip address
func GetIPAddress(c *echo.Context) string {
	return c.Request().RemoteAddr
}
