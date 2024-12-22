package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCartQuery(ctx context.Context, sc models.Scope, opt cart.CartFilter) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.buildCartQuery.BuildScopeQuery", err)
		return bson.M{}, err
	}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.IDs)}
	}

	if len(opt.ShopIDs) > 0 {
		filter["shop_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.ShopIDs)}
	}

	return filter, nil
}

func (repo implRepo) buildCartDetailQuery(ctx context.Context, sc models.Scope, id primitive.ObjectID) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.buildCartDetailQuery.BuildScopeQuery", err)
		return bson.M{}, err
	}

	filter = mongo.BuildQueryWithSoftDelete(filter)

	filter["_id"] = id

	return filter, nil
}
