package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildProductDetailQuery(ctx context.Context, ID primitive.ObjectID) (bson.M, error) {
	filter := bson.M{}
	filter["_id"] = ID

	return filter, nil
}
