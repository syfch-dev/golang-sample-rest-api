package usecase

import (
	"context"
	"my-go-project/config"
	"my-go-project/internal/models"
	"my-go-project/internal/session"

	"github.com/opentracing/opentracing-go"
)

// Session use case
type sessionUC struct {
	sessionRepo session.SessionRepository
	cfg         *config.Config
}

func NewSessionUsecase(sessionRepo session.SessionRepository, cfg *config.Config) session.UCSession {
	return &sessionUC{
		sessionRepo: sessionRepo,
		cfg:         cfg,
	}
}

// Create new session
func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.CreateSession")
	defer span.Finish()

	return u.sessionRepo.CreateSession(ctx, session, expire)
}

// Delete session by id
func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.DeleteByID")
	defer span.Finish()

	return u.sessionRepo.DeleteByID(ctx, sessionID)
}

// get session by id
func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.GetSessionByID")
	defer span.Finish()

	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}
