package shop

import "github.com/pt010104/api-golang/internal/models"

type CreateProductInput struct {
	Name            string
	Price           float32
	InventoryID     string
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
	CategoryID      []string
}
type DetailProductOutput struct {
	ID           string
	Name         string
	CategoryName []string
	Category     []models.Category

	Shop models.Shop

	Inventory models.Inventory
	Price     float32
}
type ListProductInput struct {
	CateIDs    []string
	ProductIDs []string
	ShopID     string
}
type ProductOutPutItem struct {
	P     models.Product
	Inven string
	Cate  []models.Category
}
type ListProductOutput struct {
	Products []ProductOutPutItem
	Shop     models.Shop
}
