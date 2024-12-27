package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type CreateProductInput struct {
	Name            string
	Price           float64
	InventoryID     string
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
	CategoryIDs     []string
	MediaIDs        []string
}
type DetailProductOutput struct {
	ID           string
	Name         string
	CategoryName []string
	Category     []models.Category
	Medias       []models.Media
	Shop         models.Shop
	Inventory    models.Inventory
	Price        float64
}
type GetProductOption struct {
	GetProductFilter
	PagQuery paginator.PaginatorQuery
}
type ListProductInput struct {
	GetProductFilter
}
type ProductOutPutItem struct {
	P         models.Product
	Inventory models.Inventory
	Cate      []models.Category
	Images    []models.Media
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
type UpdateProductInput struct {
	Name            string
	ID              string
	Price           float64
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
	CategoryID      []string
	MediaIDs        []string
}
type GetAllProductItem struct {
	P         models.Product
	Inventory models.Inventory
	Cate      []models.Category
	Images    []models.Media
	Shop      models.Shop
}
type GetAllProductOutput struct {
	Products []GetAllProductItem
	Pag      paginator.Paginator
}
