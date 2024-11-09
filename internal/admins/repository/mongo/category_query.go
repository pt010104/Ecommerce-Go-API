package mongo

import (
	"github.com/pt010104/api-golang/internal/admins"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildCategoryQuery(opt admins.GetCategoriesFilter) bson.M {
	filter := bson.M{}

	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": opt.IDs}
	}

	if opt.Name != "" {
		filter["name"] = opt.Name
	}
	return filter
}
