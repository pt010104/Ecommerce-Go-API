package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductOption struct {
	Name        string
	Price       float32
	InventoryID primitive.ObjectID
	ShopID      primitive.ObjectID
	CategoryID  []primitive.ObjectID
	Alias       string
	MediaIDs    []primitive.ObjectID
}
type GetProductFilter struct {
	CateIDs []string
	IDs     []string
	Search  string
	ShopID  string
}
type UpdateProductOption struct {
	Name            string
	ID              primitive.ObjectID
	Price           float32
	Model           models.Product
	Alias           string
	StockLevel      uint
	ReorderLevel    uint
	ReorderQuantity uint
	CategoryID      []primitive.ObjectID
	MediaIDs        []primitive.ObjectID
}

type UpdateManyProductsOption struct {
	ProductIDs    []primitive.ObjectID
	ReservedStock []uint
}
