package media

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadInput struct {
	Files [][]byte
}

type GetFilter struct {
	IDs    []string
	Status string
	UserID primitive.ObjectID
	ShopID primitive.ObjectID
}

type ListInput struct {
	GetFilter
}
