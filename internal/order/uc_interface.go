package order

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CheckoutUC
	OrderUC
}

type CheckoutUC interface {
	CreateCheckout(ctx context.Context, sc models.Scope, productIDs []string) (CreateCheckoutOutput, error)
}

type OrderUC interface {
	CreateOrder(ctx context.Context, sc models.Scope, input CreateOrderInput) (models.Order, error)
	DetailOrder(ctx context.Context, sc models.Scope, orderID string) (models.Order, error)
	ConsumeOrderMsg(ctx context.Context, sc models.Scope, input ConsumeOrderMsgInput) error
	ListOrder(ctx context.Context, sc models.Scope, input ListOrderInput) (ListOrderOutput, error)
}
