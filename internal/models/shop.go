package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shop struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Name   string             `bson:"name"`
	//TODO (cai nay lam sau khi co domain media): Avarta shop
	Alias      string               `bson:"alias"`
	City       string               `bson:"city"`
	Street     string               `bson:"street"`
	District   string               `bson:"district"`
	Phone      string               `bson:"phone,omitempty"`
	Followers  []primitive.ObjectID `bson:"followers,omitempty"`
	AvgRate    float64              `bson:"avg_rate"`
	UpdatedAt  time.Time            `bson:"updated_at"`
	CreatedAt  time.Time            `bson:"created_at"`
	IsVerified bool                 `bson:"is_verified"`
}
