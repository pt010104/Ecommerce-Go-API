package user

import (
	"github.com/pt010104/api-golang/internal/models"
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

type CreateKeyTokenOption struct {
	UserID       primitive.ObjectID
	SessionID    string
	RefrestToken string
	SecretKey    string
}

type UpdateUserOption struct {
	Model    models.User
	ID       string
	Email    string
	Password string
	UserName string
}

type UpdateRequestTokenOption struct {
	IsUsed *bool
	Token  string
}
