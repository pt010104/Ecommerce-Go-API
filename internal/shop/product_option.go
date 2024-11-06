package shop

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateProductOption struct {
	Name        string
	Price       float32
	InventoryID primitive.ObjectID
	ShopID      primitive.ObjectID
}
