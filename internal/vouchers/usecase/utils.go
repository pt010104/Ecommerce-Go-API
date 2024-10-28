package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pt010104/api-golang/internal/vouchers"
)

func (uc implUsecase) validateCreateVoucher(ctx context.Context, input vouchers.CreateVoucherInput) error {
	if input.ValidFrom.After(input.ValidTo) {
		return vouchers.ErrInvalidInput
	}

	if input.DiscountAmount < 0 && input.MaxDiscountAmount < 0 && input.MinimumOrderAmount < 0 {
		return vouchers.ErrInvalidInput
	}

	v, err := uc.repo.DetailVoucher(ctx, models.Scope{}, input.Code)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			uc.l.Errorf(ctx, "voucher.usecase.CreateVoucher.validateCreateVoucher.DetailVoucher: %v", err)
			return err
		}
	}
	if v.ID != primitive.NilObjectID {
		return vouchers.ErrExistCode
	}

	// if len(input.ApplicableProductIDs) > 0 {
	// 	producsts, err := uc.shopUc.ListProducts(ctx, input.ApplicableProductIDs)
	// }

	// if len(input.ApplicableCategorieIDs) > 0 {
	// 	categories, err := uc.shopUc.ListCategories(ctx, input.ApplicableCategorieIDs)
	// }

	if input.DiscountType != models.DiscountTypePercent && input.DiscountType != models.DiscountTypeFixed {
		return vouchers.ErrInvalidInput
	}

	if input.DiscountType == models.DiscountTypePercent && (input.DiscountAmount < 0 || input.DiscountAmount > 100) {
		return vouchers.ErrInvalidInput
	}

	return nil
}
