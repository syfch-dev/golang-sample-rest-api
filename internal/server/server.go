package server

import (
	"my-go-project/config"
	"my-go-project/pkg/logger"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	logger logger.Logger
}

func NewServer(cfg *config.Config, psqlDB *sqlx.DB, redisClient *redis.Client, logger logger.Logger) *Server {
	return &Server{echo: echo.New(),
		cfg:    cfg,
		logger: logger}
}

func (s *Server) Run() error {
	if s.cfg.Server.SSL {

	}

	return nil
}
