package auth

import (
	"context"
	"my-go-project/internal/models"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
