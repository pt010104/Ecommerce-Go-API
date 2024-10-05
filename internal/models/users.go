package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserName string             `bson:"user_name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}
type KeyToken struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	SecretKey    string             `bson:"secret_key"`
	RefreshToken string             `bson:"refresh_token"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
