package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type CreateProductInput struct {
	Name            string
	Price           float32
	InventoryID     string
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
	CategoryID      []string
	MediaID         string
}
type DetailProductOutput struct {
	ID           string
	Name         string
	CategoryName []string
	Category     []models.Category
	MediaID      string
	URL          string
	Shop         models.Shop

	Inventory models.Inventory
	Price     float32
}
type GetProductOption struct {
	GetProductFilter
	PagQuery paginator.PaginatorQuery
}
type ListProductInput struct {
	CateIDs    []string
	ProductIDs []string
	ShopID     string
}
type ProductOutPutItem struct {
	P       models.Product
	Inven   string
	Cate    []models.Category
	MediaID string
	URL     string
}
type ListProductOutput struct {
	Products []ProductOutPutItem
	Shop     models.Shop
}
type GetProductOutput struct {
	Products []ProductOutPutItem
	Pag      paginator.Paginator
	Shop     models.Shop
}
