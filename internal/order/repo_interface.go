package order

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repo interface {
	CreateCheckout(ctx context.Context, sc models.Scope, opt CreateCheckoutOption) (models.Checkout, error)
}
