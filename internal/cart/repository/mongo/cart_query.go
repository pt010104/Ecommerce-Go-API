package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCartQuery(sc models.Scope, opt cart.GetCartFilter, ctx context.Context) (bson.M, error) {
	filter := bson.M{}
	filter["user_id"] = bson.M{"eq": mongo.ObjectIDFromHexOrNil(sc.UserID)}
	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.IDs)}
	}
	if len(opt.ShopIDs) > 0 {
		filter["shop_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.ShopIDs)}
	}

	return filter, nil
}
func (repo implRepo) buildCartDetailQuery(ctx context.Context, ID primitive.ObjectID) (bson.M, error) {
	filter := bson.M{}
	filter = mongo.BuildQueryWithSoftDelete(filter)
	filter["_id"] = ID

	return filter, nil
}
