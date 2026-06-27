package auth

import (
	"context"
	"my-go-project/internal/models"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
}
