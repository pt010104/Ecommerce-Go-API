package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID       primitive.ObjectID `bson:"_id"`
	Street   string             `bson:"street"`
	District string             `bson:"district"`
	City     string             `bson:"city"`
	Province string             `bson:"province"`
	Phone    string             `bson:"phone"`
	Default  bool               `bson:"default"`
}

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"user_name"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	IsVerified bool               `bson:"is_verified"`
	Role       int                `bson:"role"`
	MediaID    primitive.ObjectID `bson:"media_id"`
	Address    []Address          `bson:"address"`
}
