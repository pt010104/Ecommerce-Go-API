package media

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConsumeUploadMsgInput struct {
	UserID     primitive.ObjectID
	ShopID     primitive.ObjectID
	FolderName string
	ID         string
	File       []byte
}
