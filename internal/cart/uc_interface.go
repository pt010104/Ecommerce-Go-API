package cart

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Create(sc models.Scope, ctx context.Context, input CreateCartInput, inputItem CreateCartItemInput) (models.Cart, error)
	Update(ctx context.Context, opt UpdateCartOption) (models.Cart, error)
	ListCart(sc models.Scope, ctx context.Context, opt GetCartFilter) ([]models.Cart, error)
}
