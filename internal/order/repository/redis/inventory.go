package redis

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	versionInventoryPrefix = "version_inventory"
)

func (r implRedis) GetVersionInventory(ctx context.Context, inventoryID primitive.ObjectID) (int, error) {
	key := r.buildVersionInventoryKey(inventoryID)
	val, err := r.redis.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	intVal := 0
	_, err = fmt.Sscanf(string(val), "%d", &intVal)
	if err != nil {
		return 0, err
	}

	return intVal, nil
}

func (r implRedis) IncrementVersionInventory(ctx context.Context, inventoryID primitive.ObjectID) (int, error) {
	key := r.buildVersionInventoryKey(inventoryID)
	val, err := r.redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
