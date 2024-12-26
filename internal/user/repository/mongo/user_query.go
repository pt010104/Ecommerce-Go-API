package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildUserDetailQuery(ctx context.Context, id string) (bson.M, error) {
	filter := bson.M{}
	var err error

	filter = mongo.BuildQueryWithSoftDelete(filter)

	filter["_id"], err = primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.buildUserDetailQuery: %v", err)
		return nil, err
	}

	return filter, nil
}

func (repo implRepo) buidUserQuery(ctx context.Context, opt user.GetFilter) (bson.M, error) {
	filter := bson.M{}
	var err error

	filter = mongo.BuildQueryWithSoftDelete(filter)

	if opt.ID != "" {
		filter["_id"], err = primitive.ObjectIDFromHex(opt.ID)
		if err != nil {
			repo.l.Errorf(ctx, "user.repository.mongo.buidUserQuery: %v", err)
			return nil, err
		}
	}

	if opt.Email != "" {
		filter["email"] = opt.Email
	}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": opt.IDs}
	}

	return filter, nil
}
