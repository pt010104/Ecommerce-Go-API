package shop

import ()

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
	ShopName      string
	InventoryName string
	Price         float32
}
