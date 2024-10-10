package mongo

import (
	"context"

	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildKeyTokenDetailQuery(ctx context.Context, userID string, sessionID string) (bson.M, error) {
	filter := bson.M{}
	var err error

	filter = mongo.BuildQueryWithSoftDelete(filter)

	filter["user_id"], err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.buildKeyTokenDetailQuery: %v", err)
		return nil, err
	}

	filter["session_id"] = sessionID

	return filter, nil
}
