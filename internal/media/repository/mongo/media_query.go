package mongo

import (
	"context"

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

	filter["_id"] = id

	return filter, nil
}
