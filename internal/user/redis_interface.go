package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Redis interface {
	SetSecretKey(ctx context.Context, sessionID string, secretKey string) error
	GetSecretKey(ctx context.Context, sessionID string) ([]byte, error)

	StoreSecretKey(sc models.Scope, secretKey string, ctx context.Context) error
}
