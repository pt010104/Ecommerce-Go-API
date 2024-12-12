package media

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadInput struct {
	Files [][]byte
}

type GetFilter struct {
	UserID primitive.ObjectID
	ShopID primitive.ObjectID
}
