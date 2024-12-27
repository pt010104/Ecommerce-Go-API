package usecase

import (
	"context"
	"sync"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUseCase) CreateOrder(ctx context.Context, sc models.Scope, input order.CreateOrderInput) (models.Order, error) {
	checkoutModel, err := uc.repo.DetailCheckout(ctx, sc, input.CheckoutID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.DetailCheckout", err)
		return models.Order{}, err
	}

	if checkoutModel.Status != models.CheckoutStatusPending {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.DetailCheckout", order.ErrCheckoutStatusInvalid)
		return models.Order{}, order.ErrCheckoutStatusInvalid
	}

	if checkoutModel.ExpiredAt.Before(util.Now()) {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.DetailCheckout", order.ErrCheckoutExpired)
		return models.Order{}, order.ErrCheckoutExpired
	}

	userModel, err := uc.userUC.GetModel(ctx, sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.userUC.GetModel", err)
		return models.Order{}, err
	}

	var existAddress bool
	for _, address := range userModel.Address {
		if address.ID.Hex() == input.AddressID {
			existAddress = true
			break
		}
	}

	if !existAddress {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.userUC.GetModel", order.ErrAddressNotFound)
		return models.Order{}, order.ErrAddressNotFound
	}

	var wg sync.WaitGroup
	var wgErr error
	var orderModel models.Order

	wg.Add(1)
	go func() {
		defer wg.Done()
		checkoutModel, err = uc.repo.UpdateCheckout(ctx, sc, order.UpdateCheckoutOption{
			Model:  checkoutModel,
			Status: models.CheckoutStatusConfirmed,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.UpdateCheckout", err)
			wgErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		orderModel, err = uc.repo.CreateOrder(ctx, sc, order.CreateOrderOption{
			CheckoutID:    checkoutModel.ID,
			Products:      checkoutModel.Products,
			PaymentMethod: input.PaymentMethod,
			AddressID:     input.AddressID,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.CreateOrder", err)
			wgErr = err
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return models.Order{}, wgErr
	}

	err = uc.PublishOrder(ctx, sc, rabbitmq.OrderMessage{
		OrderID:    orderModel.ID.Hex(),
		CheckoutID: checkoutModel.ID.Hex(),
		UserID:     sc.UserID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.PublishOrder", err)
		return models.Order{}, err
	}

	return orderModel, nil
}

func (uc implUseCase) DetailOrder(ctx context.Context, sc models.Scope, orderID string) (models.Order, error) {
	orderModel, err := uc.repo.DetailOrder(ctx, sc, orderID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.DetailOrder.repo.DetailOrder", err)
		return models.Order{}, err
	}

	return orderModel, nil
}

func (uc implUseCase) ListOrder(ctx context.Context, sc models.Scope, input order.ListOrderInput) ([]models.Order, error) {
	orderModel, err := uc.repo.ListOrder(ctx, sc, order.ListOrderOption{
		Status: input.Status,
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.ListOrder.repo.ListOrder", err)
		return []models.Order{}, err
	}

	return orderModel, nil
}
