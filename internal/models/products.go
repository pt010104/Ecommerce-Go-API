package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Price       float32            `bson:"price"`
	ShopID      primitive.ObjectID `bson:"shop_id"`
	InventoryID primitive.ObjectID `bson:"inventory_id"`
	CategoryID  primitive.ObjectID `bson:category_id`
	Name        string             `bson:"name"`
	Alias       string             `bson:"alias"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
