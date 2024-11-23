package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) CreateVoucher(ctx context.Context, sc models.Scope, input vouchers.CreateVoucherInput) (models.Voucher, error) {
	role := sc.Role

	util.PrintJson(input)

	if err := uc.validateCreateVoucher(ctx, input); err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.validateCreateVoucher: %v", err)
		return models.Voucher{}, err
	}

	opt := vouchers.CreateVoucherOption{
		Name:                   input.Name,
		Code:                   input.Code,
		Description:            input.Description,
		ValidFrom:              input.ValidFrom,
		ValidTo:                input.ValidTo,
		DiscountType:           input.DiscountType,
		DiscountAmount:         input.DiscountAmount,
		MaxDiscountAmount:      input.MaxDiscountAmount,
		UsageLimit:             input.UsageLimit,
		MinimumOrderAmount:     input.MinimumOrderAmount,
		ApplicableProductIDs:   input.ApplicableProductIDs,
		ApplicableCategorieIDs: input.ApplicableCategorieIDs,
		CreatedBy:              sc.UserID,
	}

	opt.ShopIDs = []string{sc.UserID}
	opt.Scope = models.ScopeShop

	if role == models.RoleAdmin {
		if len(input.ShopIDs) > 0 {
			isVerified := true
			s, err := uc.shopUc.ListShop(ctx, models.Scope{},
				shop.GetShopsFilter{
					IDs:        input.ShopIDs,
					IsVerified: &isVerified,
				},
			)
			if err != nil {
				uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.ListShop: %v", vouchers.ErrShopNotFound)
				return models.Voucher{}, vouchers.ErrShopNotFound
			}

			if len(s) != len(input.ShopIDs) {
				uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.ListShop: %v", vouchers.ErrShopNotFound)
				return models.Voucher{}, vouchers.ErrShopNotFound
			}

			opt.ShopIDs = input.ShopIDs
		} else {
			opt.Scope = models.ScopeAll
		}
	}

	v, err := uc.repo.CreateVoucher(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.CreateVoucher: %v", err)
		return models.Voucher{}, err
	}

	return v, nil
}
