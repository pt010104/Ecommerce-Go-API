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
	VoucherID     string
}

type CreateOrderOutput struct {
	OrderID string
}

type ConsumeOrderMsgInput struct {
	OrderID    string
	CheckoutID string
	UserID     string
}

type ListOrderInput struct {
	Status string
}

type ProductItem struct {
	ProductID   string
	ProductName string
	ImageURL    string
	Price       float64
	Quantity    int
}

type OrderItem struct {
	Order      models.Order
	Products   []ProductItem
	TotalPrice float64
}

type ListOrderOutput struct {
	Orders []OrderItem
}

type ListOrderShopInput struct {
	Status string
}

type ListOrderShopOutput struct {
	Orders []OrderItem
}

type UpdateOrderInput struct {
	OrderID string
	Status  string
}
