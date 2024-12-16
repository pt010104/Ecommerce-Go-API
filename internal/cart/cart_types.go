package cart

import (
	"github.com/pt010104/api-golang/internal/models"
)

type CreateCartInput struct {
	ProductID string
	Quantity  int
}

type CartItemInput struct {
	ProductID string
	Quantity  int
}

type CreateCartOutput struct {
	cart models.Cart
}
type CreateCartItemOutput struct {
	ProductID string
	Quantity  int
}

type UpdateInput struct {
	NewItemList []CartItemInput
}

type UpdateOutput struct {
	Carts []models.Cart
	Shops []models.Shop
}
