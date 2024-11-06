package shop

import (
	"github.com/pt010104/api-golang/internal/models"
)

type CreateInventoryOption struct {
	StockLevel      uint
	ReorderLevel    *uint
	ReorderQuantity *uint
}

type UpdateInventoryOption struct {
	Model           models.Inventory
	StockLevel      *uint
	ReorderLevel    *uint
	ReorderQuantity *uint
}
