package mongo

import (
	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildCategoryQuery(opt admins.GetCategoriesFilter) bson.M {
	filter := bson.M{}
	if len(opt.IDs) == 0 || len(opt.Name) == 0 {

		return filter

	}
	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.IDs)}
	}

	if opt.Name != "" {
		filter["name"] = opt.Name
	}
	return filter
}
