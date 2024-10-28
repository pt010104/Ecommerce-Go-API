package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	voucherCollection = "vouchers"
)

func (repo implRepo) getVoucherCollection() mongo.Collection {
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

	var voucher models.Voucher
	err := col.FindOne(ctx, filter).Decode(&voucher)
	if err != nil {
		repo.l.Errorf(ctx, "voucher.repo.DetailVoucher.FindOne:", err)
		return models.Voucher{}, err
	}

	return voucher, nil
}
