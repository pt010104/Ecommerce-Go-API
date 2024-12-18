package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"

	mongo1 "go.mongodb.org/mongo-driver/mongo"
)

const (
	voucherCollection = "vouchers"
)

func (repo implRepo) getVoucherCollection() mongo1.Collection {
	return *repo.database.Collection(voucherCollection)
}

func (repo implRepo) CreateVoucher(ctx context.Context, sc models.Scope, opt vouchers.CreateVoucherOption) (models.Voucher, error) {
	col := repo.getVoucherCollection()
	voucher := repo.buildVoucherModel(ctx, sc, opt)
	_, err := col.InsertOne(ctx, voucher)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repo.CreateVoucher.Insertone:", err)
		return models.Voucher{}, err
	}

	return voucher, nil
}

func (repo implRepo) DetailVoucher(ctx context.Context, sc models.Scope, code string) (models.Voucher, error) {
	col := repo.getVoucherCollection()

	filter := repo.buildVoucherDetailQuery(ctx, sc, code)
	fmt.Print(filter)
	var voucher models.Voucher
	err := col.FindOne(ctx, filter).Decode(&voucher)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repo.DetailVoucher.FindOne:", err)
		return models.Voucher{}, err
	}

	return voucher, nil
}

func (repo implRepo) ListVoucher(ctx context.Context, sc models.Scope, opt vouchers.GetVoucherFilter) ([]models.Voucher, error) {
	col := repo.getVoucherCollection()

	filter, err := repo.buildVoucherQuery(sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repository.mongo.buildVoucherQuery: %v", err)
		return []models.Voucher{}, err
	}
	fmt.Print(filter)

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repository.mongo.ListVoucer.Find: %v", err)
		return []models.Voucher{}, err
	}
	defer cursor.Close(ctx)

	var vouchers []models.Voucher
	err = cursor.All(ctx, &vouchers)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repository.mongo.ListVoucher.All: %v", err)
		return []models.Voucher{}, err
	}

	return vouchers, nil
}
func (repo implRepo) UpdateVoucher(ctx context.Context, sc models.Scope, option vouchers.UpdateVoucherOption) (models.Voucher, error) {
	col := repo.getVoucherCollection()
	option.Model.ID = mongo.ObjectIDFromHexOrNil(option.ID)
	filter := repo.buildVoucherDetailQuery(ctx, sc, option.Model.ID.Hex())

	updateData := bson.M{}

	if option.Name != "" {
		updateData["name"] = option.Name
		option.Model.Name = option.Name
	}
	//if len(option.ShopIDs) > 0 {
	//	updateData["shop_ids"] = option.ShopIDs
	//	option.Model. = option.ShopIDs
	//}
	if option.Description != "" {
		updateData["description"] = option.Description
		option.Model.Description = option.Description
	}
	if option.Code != "" {
		updateData["code"] = option.Code
		option.Model.Code = option.Code
	}
	if !option.ValidFrom.IsZero() {
		updateData["valid_from"] = option.ValidFrom
		option.Model.ValidFrom = option.ValidFrom
	}
	if !option.ValidTo.IsZero() {
		updateData["valid_to"] = option.ValidTo
		option.Model.ValidTo = option.ValidTo
	}
	if option.UsageLimit > 0 {
		updateData["usage_limit"] = option.UsageLimit
		option.Model.UsageLimit = option.UsageLimit
	}
	if len(option.ApplicableProductIDs) > 0 {
		updateData["applicable_product_ids"] = option.ApplicableProductIDs
		option.Model.ApplicableProductIDs = mongo.ObjectIDsFromHexOrNil(option.ApplicableProductIDs)
	}
	if len(option.ApplicableCategorieIDs) > 0 {
		updateData["applicable_category_ids"] = option.ApplicableCategorieIDs
		option.Model.ApplicableCategorieIDs = mongo.ObjectIDsFromHexOrNil(option.ApplicableCategorieIDs)
	}
	if option.MinimumOrderAmount > 0 {
		updateData["minimum_order_amount"] = option.MinimumOrderAmount
		option.Model.MinimumOrderAmount = option.MinimumOrderAmount
	}
	if option.DiscountType != "" {
		updateData["discount_type"] = option.DiscountType
		option.Model.DiscountType = option.DiscountType
	}
	if option.DiscountAmount > 0 {
		updateData["discount_amount"] = option.DiscountAmount
		option.Model.DiscountAmount = option.DiscountAmount
	}
	if option.MaxDiscountAmount > 0 {
		updateData["max_discount_amount"] = option.MaxDiscountAmount
		option.Model.MaxDiscountAmount = option.MaxDiscountAmount
	}

	if option.UsedCount > 0 {
		updateData["used_count"] = option.UsedCount
		option.Model.UsedCount = int(option.UsedCount)
	}

	updateData["updated_at"] = time.Now()

	update := bson.M{}
	if len(updateData) > 0 {
		update["$set"] = updateData
	}

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repo.Update.UpdateOne: %v", err)
		return models.Voucher{}, err
	}

	return option.Model, nil
}
