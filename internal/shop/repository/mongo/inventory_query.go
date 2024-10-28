package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildInventoryDetailQuery(ctx context.Context, sc models.Scope, productID primitive.ObjectID) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.buildInventoryDetailQuery.BuildScopeQuery: %v", err)
		return bson.M{}, err
	}

	filter["product_id"] = productID

	return filter, nil
}

func (repo implRepo) buildInventoryQuery(ctx context.Context, sc models.Scope, productIDs []primitive.ObjectID) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.buildGetInventoryQuery.BuildScopeQuery: %v", err)
		return bson.M{}, err
	}

	if len(productIDs) > 0 {
		filter["product_id"] = bson.M{"$in": productIDs}
	}

	return filter, nil
}
