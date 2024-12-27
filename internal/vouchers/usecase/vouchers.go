package usecase

import (
	"context"
	"sort"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) CreateVoucher(ctx context.Context, sc models.Scope, input vouchers.CreateVoucherInput) (models.Voucher, error) {
	if err := uc.validateCreateVoucher(ctx, input); err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.validateCreateVoucher: %v", err)
		return models.Voucher{}, err
	}

	opt := vouchers.CreateVoucherOption{
		Data: vouchers.VoucherData{
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
		},
	}

	opt.Data.ShopIDs = []string{sc.ShopID}
	opt.Data.Scope = models.ScopeShop

	v, err := uc.repo.CreateVoucher(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.CreateVoucher.CreateVoucher: %v", err)
		return models.Voucher{}, err
	}

	return v, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, input vouchers.DetailVoucherInput) (models.Voucher, error) {
	opt := vouchers.DetailVoucherOption{
		ID:   input.ID,
		Code: input.Code,
	}

	v, err := uc.repo.DetailVoucher(ctx, sc, opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Voucher{}, vouchers.ErrVoucherNotFound
		}
		uc.l.Errorf(ctx, "vouchers.usecase.Detail.DetailVoucher: %v", err)
		return models.Voucher{}, err

	}
	return v, nil
}

func (uc implUsecase) List(ctx context.Context, sc models.Scope, opt vouchers.GetVoucherFilter) ([]models.Voucher, error) {
	vouchers1, err := uc.repo.ListVoucher(ctx, sc, opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []models.Voucher{}, vouchers.ErrVoucherNotFound
		}

		uc.l.Errorf(ctx, "vouchers.usecase.List.ListVoucher: %v", err)
		return []models.Voucher{}, err
	}

	sort.Slice(vouchers1, func(i, j int) bool {
		return vouchers1[i].CreatedAt.After(vouchers1[j].CreatedAt)
	})

	return vouchers1, nil
}

func (uc implUsecase) ApplyVoucher(ctx context.Context, sc models.Scope, input vouchers.ApplyVoucherInput) (models.Voucher, float64, float64, error) {
	voucher, err := uc.repo.DetailVoucher(ctx, sc, vouchers.DetailVoucherOption{
		ID:   input.ID,
		Code: input.Code,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Voucher{}, 0, 0, vouchers.ErrVoucherNotFound
		}
		uc.l.Errorf(ctx, "vouchers.usecase.ApplyVoucher.DetailVoucher: %v", err)
		return models.Voucher{}, 0, 0, err
	}

	if voucher.ValidFrom.After(util.Now()) || voucher.ValidTo.Before(util.Now()) {
		return models.Voucher{}, 0, 0, vouchers.ErrVoucherExpired
	}
	if voucher.UsageLimit > 0 && uint(voucher.UsedCount) >= voucher.UsageLimit {
		return models.Voucher{}, 0, 0, vouchers.ErrVoucherExpired
	}
	if voucher.MinimumOrderAmount > 0 {
		if input.OrderAmount < voucher.MinimumOrderAmount {
			return models.Voucher{}, 0, 0, vouchers.ErrVoucherMinimumOrderAmount
		}
	}

	voucher.UsedCount++
	nv, err := uc.repo.UpdateVoucher(ctx, sc, vouchers.UpdateVoucherOption{
		Model: voucher,
	})
	if err != nil {
		uc.l.Errorf(ctx, "vouchers.usecase.ApplyVoucher.UpdateVoucher: %v", err)
		return models.Voucher{}, 0, 0, err
	}

	discountAmount := 0.0
	if voucher.DiscountType == models.DiscountTypeFixed {
		discountAmount = voucher.DiscountAmount
		if voucher.MaxDiscountAmount > 0 && discountAmount > voucher.MaxDiscountAmount {
			discountAmount = voucher.MaxDiscountAmount
		}
	} else if voucher.DiscountType == models.DiscountTypePercent {
		discountAmount = input.OrderAmount * voucher.DiscountAmount / 100
		if voucher.MaxDiscountAmount > 0 && discountAmount > voucher.MaxDiscountAmount {
			discountAmount = voucher.MaxDiscountAmount
		}
	}

	orderAmount := input.OrderAmount - discountAmount
	if orderAmount < 0 {
		orderAmount = 0
	}

	return nv, orderAmount, discountAmount, nil
}
