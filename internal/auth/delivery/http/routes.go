package auth

import (
	"my-go-project/internal/auth"

	"github.com/labstack/echo/v5"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers) {
	authGroup.POST("/register", h.Register())
}