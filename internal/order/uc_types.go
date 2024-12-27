package order

import (
	"time"

	"github.com/pt010104/api-golang/internal/models"
)

type CreateCheckoutOutput struct {
	CheckoutID       string
	ExpiredAt        time.Time
	TotalPriceByShop map[string]float64
	TotalPrice       float64
	Products         []models.Product
	QuantityMap      map[string]int
	Shops            []models.Shop
	ImageURLs        []string
}

// ------------------------------ //

type CreateOrderInput struct {
	CheckoutID    string
	PaymentMethod string
	AddressID     string
}

type CreateOrderOutput struct {
	OrderID string
}

type ConsumeOrderMsgInput struct {
	OrderID    string
	CheckoutID string
	UserID     string
}
