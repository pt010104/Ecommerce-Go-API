package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type address struct {
	Street   string `bson:"street"`
	District string `bson:"district"`
	City     string `bson:"city"`
}

type Shop struct {
	ID        int                  `bson:"_id"`
	UserID    int                  `bson:"user_id"`
	Name      string               `bson:"name"`
	Alias     string               `bson:"alias"`
	Address   address              `bson:"address"`
	Phone     string               `bson:"phone,omitempty"`
	Followers []primitive.ObjectID `bson:"followers,omitempty"`
	AvgRate   float64              `bson:"avg_rate"`
	UpdatedAt string               `bson:"updated_at"`
	CreatedAt string               `bson:"created_at"`
}
