package media

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConsumeUploadMsgInput struct {
	ID         primitive.ObjectID
	UserID     primitive.ObjectID
	ShopID     primitive.ObjectID
	FolderName string
	FileName   string
	File       []byte
}
