package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (r implRepository) buildDetailQuery(ctx context.Context, sc models.Scope, id string) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, r.l, sc)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.buildDetailQuery.BuildScopeQuery: %v", err)
		return nil, err
	}

	filter["_id"] = mongo.ObjectIDFromHexOrNil(id)

	return filter, nil
}

func (r implRepository) buildQuery(ctx context.Context, sc models.Scope, opt media.GetFilter) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, r.l, sc)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.buildQuery.BuildScopeQuery: %v", err)
		return nil, err
	}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": opt.IDs}
	}

	if opt.Status != "" {
		filter["status"] = opt.Status
	}

	return filter, nil
}
