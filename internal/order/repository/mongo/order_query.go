package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildOrderDetailQuery(ctx context.Context, sc models.Scope, orderID string) (bson.M, error) {
	filter, err := mongo.BuildScopeQuery(ctx, repo.l, sc)
	if err != nil {
		repo.l.Errorf(ctx, "Checkout.Repo.buildCheckoutDetailQuery.BuildScopeQuery", err)
		return nil, err
	}

	filter["_id"] = mongo.ObjectIDFromHexOrNil(orderID)

	return filter, nil
}
