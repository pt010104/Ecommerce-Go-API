package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeyToken struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	SecretKey    string             `bson:"secret_key"`
	RefreshToken string             `bson:"refresh_token"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
