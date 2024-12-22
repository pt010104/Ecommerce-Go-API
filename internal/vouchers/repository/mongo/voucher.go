package mongo

import (
	"context"
	"fmt"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"

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

func (repo implRepo) DetailVoucher(ctx context.Context, sc models.Scope, opt vouchers.DetailVoucherOption) (models.Voucher, error) {
	col := repo.getVoucherCollection()

	filter := repo.buildVoucherDetailQuery(ctx, sc, opt)
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

	filter := repo.buildVoucherDetailQuery(ctx, sc, vouchers.DetailVoucherOption{
		ID: option.Model.ID.Hex(),
	})

	nm, update, err := repo.buildVoucherUpdate(ctx, sc, option)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repo.Update.buildVoucherUpdate: %v", err)
		return models.Voucher{}, err
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repo.Update.UpdateOne: %v", err)
		return models.Voucher{}, err
	}

	return nm, nil
}
