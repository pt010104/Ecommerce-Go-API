package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Items     []CartItem         `bson:"items" json:"items"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	ShopID    primitive.ObjectID `bson:"shop_id"`
}

type CartItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"productId"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	AddedAt   time.Time          `bson:"addedAt" json:"addedAt"`
}
