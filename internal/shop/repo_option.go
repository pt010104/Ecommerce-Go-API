package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateShopOption struct {
	Name     string
	Alias    string
	City     string
	Street   string
	District string
	Phone    string
}

type GetOption struct {
	GetShopsFilter
	PagQuery paginator.PaginatorQuery
}

type DetailOption struct {
	ID string
}
type UpdateOption struct {
	Model      models.Shop
	Name       string
	Alias      string
	City       string
	Street     string
	District   string
	Phone      string
	IsVerified bool
}

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
