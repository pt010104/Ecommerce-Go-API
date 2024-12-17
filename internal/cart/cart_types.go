package cart

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
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
type ProductItem struct {
	ProductID   string
	Medias      []models.Media
	Quantity    int
	Price       float32
	ProductName string
}

type GetCartItem struct {
	Products []ProductItem

	CartProductMediaMap map[string][]models.Media
	Cart                models.Cart
	Shop                models.Shop
}
type GetCartOutput struct {
	CartOutPut []GetCartItem
	Paginator  paginator.Paginator
}
