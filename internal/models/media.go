package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID        primitive.ObjectID `bson:"_id"`
	URL       string             `bson:"url"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ShopID    primitive.ObjectID `bson:"shop_id,omitempty"`
	FileName  string             `bson:"file_name"`
	Folder    string             `bson:"folder"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

const (
	MediaStatusPending  = "pending"
	MediaStatusUploaded = "uploaded"
	MediaStatusFailed   = "failed"
	MediaStatusDrafted  = "drafted"
)
