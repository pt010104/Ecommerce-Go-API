package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestToken struct {
	ID        primitive.ObjectID `bson:"id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdateAt  time.Time          `bson:"update_at"`
	Is_Used   bool               `bson:"is_used"`
}
