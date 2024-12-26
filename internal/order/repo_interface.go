package order

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repo interface {
	OrderRepo
	CheckoutRepo
}

type OrderRepo interface {
	CreateOrder(ctx context.Context, sc models.Scope, opt CreateOrderOption) (models.Order, error)
	DetailOrder(ctx context.Context, sc models.Scope, orderID string) (models.Order, error)
}

type CheckoutRepo interface {
	CreateCheckout(ctx context.Context, sc models.Scope, opt CreateCheckoutOption) (models.Checkout, error)
	DetailCheckout(ctx context.Context, sc models.Scope, checkoutID string) (models.Checkout, error)
	UpdateCheckout(ctx context.Context, sc models.Scope, opt UpdateCheckoutOption) (models.Checkout, error)
}
