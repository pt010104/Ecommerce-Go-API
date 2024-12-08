package media

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadInput struct {
	UserID primitive.ObjectID
	ShopID primitive.ObjectID
	Files  [][]byte
}

type GetFilter struct {
	UserID primitive.ObjectID
	ShopID primitive.ObjectID
}
