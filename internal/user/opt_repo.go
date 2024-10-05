package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepoOption struct {
	UserName string
	Password string
	Email    string
}

type KeyTokenRepoOption struct {
	UserID primitive.ObjectID
}
