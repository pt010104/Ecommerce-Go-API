package order

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateCheckoutOption struct {
	ProductIDs []primitive.ObjectID
}
