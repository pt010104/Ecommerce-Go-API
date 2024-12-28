package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildProductDetailQuery(ctx context.Context, ID primitive.ObjectID) (bson.M, error) {
	filter := bson.M{}
	filter = mongo.BuildQueryWithSoftDelete(filter)
	filter["_id"] = ID

	return filter, nil
}
func (repo implRepo) buildProductDetailUpdateQuery(ctx context.Context, ID string) (bson.M, error) {
	filter := bson.M{}
	filter = mongo.BuildQueryWithSoftDelete(filter)
	filter["_id"] = mongo.ObjectIDFromHexOrNil(ID)

	return filter, nil
}
func (repo implRepo) buildProductQuery(sc models.Scope, opt shop.GetProductFilter) (bson.M, error) {
	filter := bson.M{}

	if opt.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": opt.Search, "$options": "i"}},
			{"shop_name": bson.M{"$regex": opt.Search, "$options": "i"}},
			{"alias": bson.M{"$regex": util.BuildAlias(opt.Search), "$options": "i"}},
		}
	}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.IDs)}
	}
	if opt.ShopID != "" {
		filter["shop_id"] = bson.M{"$eq": mongo.ObjectIDFromHexOrNil(opt.ShopID)}
	}

	if len(opt.InventoryIDs) > 0 {
		filter["inventory_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.InventoryIDs)}
	}

	if len(opt.CateIDs) > 0 {
		filter["categoryid"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.CateIDs)}
	}
	return filter, nil
}
func (repo implRepo) buildProductDeleteQuery(sc models.Scope, ctx context.Context, ids []string) (bson.M, error) {
	filter := bson.M{}
	filter["shop_id"] = bson.M{"$eq": mongo.ObjectIDFromHexOrNil(sc.ShopID)}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(ids)}
	}

	return filter, nil
}
