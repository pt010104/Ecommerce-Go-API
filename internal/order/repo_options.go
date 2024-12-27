package order

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCheckoutOption struct {
	Products []models.OrderProduct
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
}
