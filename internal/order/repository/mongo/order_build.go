package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildOrderModel(ctx context.Context, sc models.Scope, opt order.CreateOrderOption) (models.Order, error) {
	now := time.Now()

	order := models.Order{
		ID:            primitive.NewObjectID(),
		UserID:        mongo.ObjectIDFromHexOrNil(sc.UserID),
		CheckoutID:    opt.CheckoutID,
		Products:      opt.Products,
		Status:        models.OrderStatusPending,
		PaymentMethod: opt.PaymentMethod,
		AddressID:     mongo.ObjectIDFromHexOrNil(opt.AddressID),
		TotalPrice:    opt.TotalPrice,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return order, nil
}
