package cart

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repo interface {
	Create(ctx context.Context, sc models.Scope, opt CreateCartOption) (models.Cart, error)
	GetOne(ctx context.Context, sc models.Scope, opt GetOneOption) (models.Cart, error)
	Update(ctx context.Context, sc models.Scope, opt UpdateCartOption) (models.Cart, error)
	ListCart(ctx context.Context, sc models.Scope, opt ListOption) ([]models.Cart, error)
}
