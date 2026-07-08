package usecase

import (
	"context"
	"fmt"
	"my-go-project/config"
	"my-go-project/internal/auth"
	"my-go-project/internal/models"
	"my-go-project/pkg/logger"
	"my-go-project/pkg/utils"
	"net/http"

	"my-go-project/pkg/httpErrors"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg       *config.Config
	logger    logger.Logger
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
}

func NewAuthUseCase(cfg *config.Config, logger logger.Logger, authRepo auth.Repository) auth.UseCase {
	return &authUC{cfg: cfg, logger: logger, authRepo: authRepo}
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

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	defer span.Finish()

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID.String()), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
	}

	user.SanitizePassword()

	return user, nil
}

func (u *authUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
