package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildVoucherModel(ctx context.Context, sc models.Scope, opt vouchers.CreateVoucherOption) models.Voucher {
	now := util.Now()

	v := models.Voucher{
		ID:                     primitive.NewObjectID(),
		Name:                   opt.Data.Name,
		Code:                   opt.Data.Code,
		Description:            opt.Data.Description,
		ValidFrom:              opt.Data.ValidFrom,
		ValidTo:                opt.Data.ValidTo,
		DiscountType:           opt.Data.DiscountType,
		DiscountAmount:         opt.Data.DiscountAmount,
		MaxDiscountAmount:      opt.Data.MaxDiscountAmount,
		UsageLimit:             opt.Data.UsageLimit,
		MinimumOrderAmount:     opt.Data.MinimumOrderAmount,
		CreatedBy:              mongo.ObjectIDFromHexOrNil(opt.Data.CreatedBy),
		ShopIDs:                mongo.ObjectIDsFromHexOrNil(opt.Data.ShopIDs),
		ApplicableCategorieIDs: mongo.ObjectIDsFromHexOrNil(opt.Data.ApplicableCategorieIDs),
		ApplicableProductIDs:   mongo.ObjectIDsFromHexOrNil(opt.Data.ApplicableProductIDs),
		Scope:                  opt.Data.Scope,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	return v
}

func (repo implRepo) buildVoucherUpdate(ctx context.Context, sc models.Scope, opt vouchers.UpdateVoucherOption) (models.Voucher, bson.M, error) {
	now := util.Now()
	updateData := bson.M{
		"updated_at":              now,
		"name":                    opt.Data.Name,
		"code":                    opt.Data.Code,
		"description":             opt.Data.Description,
		"valid_from":              opt.Data.ValidFrom,
		"valid_to":                opt.Data.ValidTo,
		"usage_limit":             opt.Data.UsageLimit,
		"applicable_product_ids":  mongo.ObjectIDsFromHexOrNil(opt.Data.ApplicableProductIDs),
		"applicable_category_ids": mongo.ObjectIDsFromHexOrNil(opt.Data.ApplicableCategorieIDs),
		"shop_ids":                mongo.ObjectIDsFromHexOrNil(opt.Data.ShopIDs),
		"minimum_order_amount":    opt.Data.MinimumOrderAmount,
		"discount_type":           opt.Data.DiscountType,
		"discount_amount":         opt.Data.DiscountAmount,
		"max_discount_amount":     opt.Data.MaxDiscountAmount,
		"scope":                   opt.Data.Scope,
	}
	opt.Model.UpdatedAt = now
	opt.Model.Name = opt.Data.Name
	opt.Model.Code = opt.Data.Code
	opt.Model.Description = opt.Data.Description
	opt.Model.ValidFrom = opt.Data.ValidFrom
	opt.Model.ValidTo = opt.Data.ValidTo
	opt.Model.UsageLimit = opt.Data.UsageLimit
	opt.Model.ApplicableProductIDs = mongo.ObjectIDsFromHexOrNil(opt.Data.ApplicableProductIDs)
	opt.Model.ApplicableCategorieIDs = mongo.ObjectIDsFromHexOrNil(opt.Data.ApplicableCategorieIDs)
	opt.Model.MinimumOrderAmount = opt.Data.MinimumOrderAmount
	opt.Model.DiscountType = opt.Data.DiscountType
	opt.Model.DiscountAmount = opt.Data.DiscountAmount
	opt.Model.MaxDiscountAmount = opt.Data.MaxDiscountAmount
	opt.Model.ShopIDs = mongo.ObjectIDsFromHexOrNil(opt.Data.ShopIDs)
	opt.Model.Scope = opt.Data.Scope

	if opt.Data.UsedCount > 0 {
		updateData["used_count"] = opt.Data.UsedCount
		opt.Model.UsedCount = opt.Data.UsedCount
	}

	update := bson.M{}
	if len(updateData) > 0 {
		update["$set"] = updateData
	}

	return opt.Model, update, nil
}
