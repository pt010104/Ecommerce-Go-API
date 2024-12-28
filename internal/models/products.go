package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID   `bson:"_id"`
	Price       float64              `bson:"price"`
	ShopID      primitive.ObjectID   `bson:"shop_id"`
	InventoryID primitive.ObjectID   `bson:"inventory_id"`
	CategoryID  []primitive.ObjectID `bson:category_id`
	Name        string               `bson:"name"`
	Description string               `bson:"description"`
	Alias       string               `bson:"alias"`
	View        int                  `bson:"view"`
	ViewTrend   int                  `bson:"view_trend"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
	MediaIDs    []primitive.ObjectID `bson:"media_ids"`
	Thumbnail   primitive.ObjectID   `bson:"thumbnail"`
}
