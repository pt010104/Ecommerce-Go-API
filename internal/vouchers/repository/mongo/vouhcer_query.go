package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"

	"github.com/pt010104/api-golang/internal/vouchers"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildVoucherDetailQuery(ctx context.Context, sc models.Scope, opt vouchers.DetailVoucherOption) bson.M {
	filter := bson.M{}

	if opt.ID != "" {
		filter["_id"] = mongo.ObjectIDFromHexOrNil(opt.ID)
	}

	if opt.Code != "" {
		filter["code"] = opt.Code
	}

	return filter

}
func (repo implRepo) buildVoucherQuery(sc models.Scope, opt vouchers.GetVoucherFilter) (bson.M, error) {
	filter := bson.M{}
	filter["scope"] = opt.Scope
	if len(opt.IDs) > 0 {
		filter["_id"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.IDs)}
	}

	if len(opt.Codes) > 0 {
		filter["code"] = bson.M{"$in": opt.Codes}
	}
	if len(opt.ApplicableCategorieIDs) > 0 {
		filter["applicable_category_ids"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.ApplicableCategorieIDs)}
	}
	if len(opt.ApplicableProductIDs) > 0 {
		filter["applicable_product_ids"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.ApplicableProductIDs)}
	}
	if len(opt.ShopIDs) > 0 {
		filter["shop_ids"] = bson.M{"$in": mongo.ObjectIDsFromHexOrNil(opt.ShopIDs)}
	}

	return filter, nil
}
