package mongo

import (
	"context"

	//"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildRequestTokenDetailQuery(ctx context.Context, token string) bson.M {
	filter := bson.M{}

	filter = mongo.BuildQueryWithSoftDelete(filter)

	filter["token"] = token

	return filter
}
