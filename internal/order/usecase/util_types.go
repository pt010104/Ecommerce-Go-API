package usecase

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

type getDataOutput struct {
	ShopIDs    []string
	Shops      []models.Shop
	CartItems  []models.CartItem
	ProductMap map[string]shop.ProductOutPutItem
}
