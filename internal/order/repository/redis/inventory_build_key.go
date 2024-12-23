package redis

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r implRedis) buildVersionInventoryKey(inventoryID primitive.ObjectID) string {
	return fmt.Sprintf("%s:%s", versionInventoryPrefix, inventoryID.Hex())
}
