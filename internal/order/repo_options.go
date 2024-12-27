package order

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCheckoutOption struct {
	Products   []models.OrderProduct
	TotalPrice float64
}

type UpdateCheckoutOption struct {
	Model  models.Checkout
	Status string
}

type CreateOrderOption struct {
	CheckoutID    primitive.ObjectID
	Products      []models.OrderProduct
	PaymentMethod string
	AddressID     string
	TotalPrice    float64
}

type ListOrderOption struct {
	Status string
}

type UpdateOrderOption struct {
	Model  models.Order
	Status string
}
