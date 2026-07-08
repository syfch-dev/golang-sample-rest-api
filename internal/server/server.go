package server

import (
	"context"
	"my-go-project/config"
	"my-go-project/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	echo        *echo.Echo
	cfg         *config.Config
	logger      logger.Logger
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewServer(cfg *config.Config, psqlDB *sqlx.DB, redisClient *redis.Client, logger logger.Logger) *Server {
	return &Server{
		echo:        echo.New(),
		cfg:         cfg,
		logger:      logger,
		db:          psqlDB,
		redisClient: redisClient,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	// pprof debug sunucusu arka planda çalışır
	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
		if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	// Ctrl+C / SIGTERM sinyalinde context iptal edilir;
	// StartConfig bunu yakalayıp graceful shutdown başlatır
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	sc := echo.StartConfig{
		Address:         s.cfg.Server.Port,
		GracefulTimeout: ctxTimeout * time.Second,
		BeforeServeFunc: func(srv *http.Server) error {
			srv.ReadTimeout = s.cfg.Server.ReadTimeout
			srv.WriteTimeout = s.cfg.Server.WriteTimeout
			srv.MaxHeaderBytes = maxHeaderBytes
			return nil
		},
	}

	s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)

	var err error
	if s.cfg.Server.SSL {
		err = sc.StartTLS(ctx, s.echo, certFile, keyFile)
	} else {
		err = sc.Start(ctx, s.echo)
	}

	s.logger.Info("Server Exited Properly")
	return err
}
