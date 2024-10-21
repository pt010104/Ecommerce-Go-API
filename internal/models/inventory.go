package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID              primitive.ObjectID `bson:"_id"`
	ProductID       primitive.ObjectID `bson:"product_id"`
	StockLevel      int                `bson:"stock_level"`
	ReservedLevel   int                `bson:"reserved_level"`
	ReorderLevel    *int               `bson:"reorder_level"`
	ReorderQuantity *int               `bson:"reorder_quantity"`
	SoldQuantity    int                `bson:"sold_quantity"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}
