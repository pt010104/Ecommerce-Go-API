package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Register(ctx context.Context, sc models.Scope, input RegisterInput) (models.Shop, error)
}
