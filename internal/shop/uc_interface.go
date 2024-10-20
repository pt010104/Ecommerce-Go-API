package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, sc models.Scope, input CreateInput) (models.Shop, error)
	Get(ctx context.Context, sc models.Scope, input GetInput) (GetOutput, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	Delete(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	Update(ctx context.Context, sc models.Scope, input UpdateInput) (models.Shop, error)
}
