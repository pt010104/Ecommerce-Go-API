package models

import (
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID     primitive.ObjectID `bson:"_id"`
	Price  float32            `bson:"price"`
	ShopID primitive.ObjectID `bson:"shop_id"`
	Avatar url.URL            `bson:"avatar"`
}
