package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
)

func (uc implUseCase) PublishOrder(ctx context.Context, sc models.Scope, input rabbitmq.OrderMessage) error {
	err := uc.prod.PublishOrderMsg(ctx, input)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.PublishOrder.prod.PublishOrderMsg: %v", err)
		return err
	}

	return nil
}
