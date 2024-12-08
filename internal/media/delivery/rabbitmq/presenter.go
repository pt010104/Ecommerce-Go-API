package rabbitmq

import "go.mongodb.org/mongo-driver/bson/primitive"

type UploadMessage struct {
	UserID     primitive.ObjectID `json:"user_id"`
	ShopID     primitive.ObjectID `json:"shop_id"`
	FileName   string             `json:"file_name"`
	File       []byte             `json:"file"`
	FolderName string             `json:"folder_name"`
}
