package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildVoucherModel(ctx context.Context, sc models.Scope, opt vouchers.CreateVoucherOption) models.Voucher {
	now := util.Now()

	v := models.Voucher{
		ID:                 primitive.NewObjectID(),
		Name:               opt.Name,
		Code:               opt.Code,
		Description:        opt.Description,
		ValidFrom:          opt.ValidFrom,
		ValidTo:            opt.ValidTo,
		DiscountType:       opt.DiscountType,
		DiscountAmount:     opt.DiscountAmount,
		MaxDiscountAmount:  opt.MaxDiscountAmount,
		UsageLimit:         opt.UsageLimit,
		MinimumOrderAmount: opt.MinimumOrderAmount,
		CreatedBy:          mongo.ObjectIDFromHexOrNil(opt.CreatedBy),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if len(opt.ApplicableCategorieIDs) > 0 {
		applicableCategorieIDs, _ := mongo.ObjectIDsFromHex(opt.ApplicableCategorieIDs)
		v.ApplicableCategorieIDs = applicableCategorieIDs
	}

	if len(opt.ApplicableProductIDs) > 0 {
		applicableProductIDs, _ := mongo.ObjectIDsFromHex(opt.ApplicableProductIDs)
		v.ApplicableProductIDs = applicableProductIDs
	}

	return v
}
