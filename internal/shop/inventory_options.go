package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInventoryOption struct {
	ProductID       primitive.ObjectID
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
