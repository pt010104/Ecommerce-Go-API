package user

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserOption struct {
	Name     string
	Password string
	Email    string
}

type GetFilter struct {
	ID    string
	IDs   []string
	Email string
}

type GetUserOption struct {
	GetFilter
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
	Model      models.User
	Email      string
	Password   string
	Name       string
	IsVerified bool
	MediaID    string
	Address    []models.Address
}
type UpdateKeyTokenInput struct {
	UserID       string
	SessionID    string
	RefreshToken string
}
type UpdateRequestTokenOption struct {
	IsUsed *bool
	Token  string
}

type AddAddressOption struct {
	Street   string
	District string
	City     string
	Province string
	Phone    string
	Default  bool
}

type ListUserOption struct {
	GetFilter
}
