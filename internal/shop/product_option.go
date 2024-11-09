package shop

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateProductOption struct {
	Name        string
	Price       float32
	InventoryID primitive.ObjectID
	ShopID      primitive.ObjectID
	CategoryID  []primitive.ObjectID
	Alias       string
}
type GetProductFilter struct {
	IDs    []string
	Search string
	ShopID []string
}
