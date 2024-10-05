package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserName string             `bson:"user_name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}
