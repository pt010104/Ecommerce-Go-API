package cart

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type CreateCartOption struct {
	ShopID       string
	CartItemList []models.CartItem
}

type UpdateCartOption struct {
	Model       models.Cart
	NewItemList []models.CartItem
}
type UpdateCartItemOption struct {
	Quantity int
}
type CartFilter struct {
	IDs     []string
	ShopIDs []string
	ID      string
}

type ListOption struct {
	CartFilter
}

type GetOneOption struct {
	CartFilter
}
type GetOption struct {
	CartFilter
	PagQuery paginator.PaginatorQuery
}
