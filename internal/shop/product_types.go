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
	ID            string
	Name          string
	CategoryName  []string
	Category      []models.Category
	ShopName      string
	Shop          models.Shop
	InventoryName string
	Inventory     models.Inventory
	Price         float32
}

type ListProductOutput struct {
	List []DetailProductOutput
}
