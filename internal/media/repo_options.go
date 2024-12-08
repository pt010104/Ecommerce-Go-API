package media

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadOption struct {
	UserID   primitive.ObjectID
	ShopID   primitive.ObjectID
	Folder   string
	FileName string
}

type UpdateOption struct {
	Model  models.Media
	Status string
	URL    string
}
