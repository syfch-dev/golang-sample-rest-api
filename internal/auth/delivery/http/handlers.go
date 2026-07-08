package http

import (
	"net/http"

	"my-go-project/config"
	"my-go-project/internal/auth"
	"my-go-project/internal/models"
	"my-go-project/pkg/httpErrors"
	"my-go-project/pkg/logger"

	"github.com/labstack/echo/v5"
	"github.com/opentracing/opentracing-go"
)

type authHandlers struct {
	cfg    *config.Config
	logger logger.Logger
	authUC auth.UseCase
}

func NewAuthHandlers(cfg *config.Config, logger logger.Logger, authUC auth.UseCase) auth.Handlers {
	return &authHandlers{cfg: cfg, logger: logger, authUC: authUC}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c *echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "auth.Register")
		defer span.Finish()

		user := &models.User{}
		if err := c.Bind(user); err != nil {
			code, resp := httpErrors.ErrorResponse(err)
			return c.JSON(code, resp)
		}

		createdUser, err := h.authUC.Register(ctx, user)
		if err != nil {
			code, resp := httpErrors.ErrorResponse(err)
			return c.JSON(code, resp)
		}

		return c.JSON(http.StatusCreated, createdUser)
	}
}
