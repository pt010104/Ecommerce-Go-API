package order

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateCheckout(ctx context.Context, sc models.Scope, productIDs []string) (CreateCheckoutOutput, error)
}
