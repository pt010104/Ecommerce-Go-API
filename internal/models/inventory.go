package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID              primitive.ObjectID `bson:"_id"`
	ProductID       primitive.ObjectID `bson:"product_id"`
	StockLevel      uint               `bson:"stock_level"`
	ReservedLevel   uint               `bson:"reserved_level"`
	ReorderLevel    *uint              `bson:"reorder_level,omitempty"`
	ReorderQuantity *uint              `bson:"reorder_quantity,omitempty"`
	SoldQuantity    uint               `bson:"sold_quantity"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}
