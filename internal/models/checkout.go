package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderProduct struct {
	ID       primitive.ObjectID `bson:"_id"`
	Quantity int                `bson:"quantity"`
}

type Checkout struct {
	ID        primitive.ObjectID `bson:"_id"`
	Products  []OrderProduct     `bson:"products"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Status    string             `bson:"status"`
	ExpiredAt time.Time          `bson:"expired_at"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

const (
	CheckoutStatusPending   = "pending"
	CheckoutStatusCanceled  = "canceled"
	CheckoutStatusConfirmed = "confirmed"
)
