package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Redis interface {
	GetSecretKey(ctx context.Context, sc models.Scope) (string, error)
	StoreSecretKey(sc models.Scope, secretKey string, ctx context.Context) error

	DetailUser(ctx context.Context, userID string) (models.User, error)
	StoreUser(ctx context.Context, user models.User) error
}
