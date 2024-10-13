package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildShopQuery(ctx context.Context, sc models.Scope, opt shop.GetOption) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.buildCandidateQuery.BuildScopeQuery: %v", err)
		return nil, err
	}

	filter = mongo.BuildQueryWithSoftDelete(filter)

	if opt.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": opt.Search, "$options": "i"}},
			{"phone": bson.M{"$regex": opt.Search, "$options": "i"}},
			{"alias": bson.M{"$regex": util.BuildAlias(opt.Search), "$options": "i"}},
		}
	}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": opt.IDs}
	}

	return filter, nil
}

func (repo implRepo) buildShopDetailQuery(ctx context.Context, sc models.Scope, id string) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.buildCandidateDetailQuery.BuildScopeQuery: %v", err)
		return nil, err
	}

	filter = mongo.BuildQueryWithSoftDelete(filter)

	if id != "" {
		filter["_id"] = id
	}

	return filter, nil
}
