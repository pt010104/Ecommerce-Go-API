package order

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Redis interface {
	GetVersionInventory(ctx context.Context, inventoryID primitive.ObjectID) (int, error)
	IncrementVersionInventory(ctx context.Context, inventoryID primitive.ObjectID) (int, error)
}
