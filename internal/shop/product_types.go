package shop

import ()

type CreateProductInput struct {
	Name            string
	Price           float32
	InventoryID     string
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
}
