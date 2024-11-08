package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildProductDetailQuery(ctx context.Context, ID primitive.ObjectID) (bson.M, error) {
	filter := bson.M{}
	filter["_id"] = ID

	return filter, nil
}

func (repo implRepo) buildProductQuery(opt shop.GetProductFilter) (bson.M, error) {
	filter := bson.M{}

	if opt.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": opt.Search, "$options": "i"}},
		}
	}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.IDs)}
	}
	return filter, nil
}
