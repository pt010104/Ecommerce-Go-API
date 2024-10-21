package inventory

import "github.com/pt010104/api-golang/internal/models"

type CreateInput struct {
	ProductID       string
	StockLevel      int
	ReorderLevel    *int
	ReorderQuantity *int
}

type CreateOutput struct {
	Inventory models.Inventory
}
