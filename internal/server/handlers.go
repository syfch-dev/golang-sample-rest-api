package server

import (
	authHttp "my-go-project/internal/auth/delivery/http"
	authRepository "my-go-project/internal/auth/repository"
	authUseCase "my-go-project/internal/auth/usecase"
	apiMiddlewares "my-go-project/internal/middleware"

	sessionRepository "my-go-project/internal/session/repository"

	sessionUseCase "my-go-project/internal/session/usecase"

	"my-go-project/pkg/csrf"
	"my-go-project/pkg/metric"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
	"github.com/swaggo/swag/example/basic/docs"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)
	_ = metrics

	authRepo := authRepository.NewAuthRepository(s.db)
	authUC := authUseCase.NewAuthUseCase(s.cfg, s.logger, authRepo)
	authHandlers := authHttp.NewAuthHandlers(s.cfg, s.logger, authUC)

	sessionRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	sessionUC := sessionUseCase.NewSessionUsecase(sessionRepo, s.cfg)

	mw := apiMiddlewares.NewMiddlewareManager(sessionUC, authUC, s.cfg, []string{"*"}, s.logger)
	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}
	docs.SwaggerInfo.Title = "Go example REST API"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RequestID())
	//e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c *echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit(2 << 20)) // 2MB

	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authHttp.MapAuthRoutes(authGroup, authHandlers)

	v1.GET("/health", func(c *echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", c.Request().Header.Get(echo.HeaderXRequestID))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
