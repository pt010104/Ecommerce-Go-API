package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserOption struct {
	UserName string
	Password string
	Email    string
}

type GetUserOption struct {
	ID    string
	Email string
}

type KeyTokenRepoOption struct {
	UserID primitive.ObjectID
}
