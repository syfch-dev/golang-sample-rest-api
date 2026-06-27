package usecase

import (
	"context"
	"my-go-project/config"
	"my-go-project/internal/auth"
	"my-go-project/internal/models"
	"my-go-project/pkg/logger"
	"net/http"

	"my-go-project/pkg/httpErrors"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg      *config.Config
	logger   logger.Logger
	authRepo auth.Repository
}

func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	user.PrepareCreate()
	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}
	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()


	return nil, nil
}
