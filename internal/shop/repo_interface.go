package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type Repo interface {
	Create(ctx context.Context, sc models.Scope, opt CreateOption) (models.Shop, error)
	Get(ctx context.Context, sc models.Scope, opt GetOption) ([]models.Shop, paginator.Paginator, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	Delete(ctx context.Context, sc models.Scope) error
	FindByid(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	Update(ctx context.Context, sc models.Scope, option UpdateOption) error
}
