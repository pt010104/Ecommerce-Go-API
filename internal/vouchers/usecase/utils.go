package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
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

	if len(input.ApplicableProductIDs) > 0 {
		productOutput, err := uc.shopUc.ListProduct(ctx, models.Scope{}, shop.ListProductInput{
			GetProductFilter: shop.GetProductFilter{
				IDs: input.ApplicableProductIDs,
			},
		})
		if err != nil {
			uc.l.Errorf(ctx, "voucher.usecase.CreateVoucher.validateCreateVoucher.ListProduct: %v", err)
			return err
		}

		if len(productOutput.Products) != len(input.ApplicableProductIDs) {
			return vouchers.ErrInvalidInput
		}
	}

	if input.DiscountType != models.DiscountTypePercent && input.DiscountType != models.DiscountTypeFixed {
		return vouchers.ErrInvalidInput
	}

	if input.DiscountType == models.DiscountTypePercent && (input.DiscountAmount < 0 || input.DiscountAmount > 100) {
		return vouchers.ErrInvalidInput
	}

	return nil
}
