package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

func (uc implUsecase) VerifyShop(ctx context.Context, sc models.Scope, input admin.VerifyShopInput) (models.Shop, error) {
	if sc.Role == 0 {
		return models.Shop{}, admin.ErrNoPermission
	}
	shop, err := uc.shopUc.Update(ctx, models.Scope{}, shop.UpdateInput{
		ShopID:     input.ShopID,
		IsVerified: true,
	})
	if err != nil {
		uc.l.Errorf(ctx, "admin.usecase.Verifyshop.shopUC.update:", err)
		return models.Shop{}, err
	}
	return shop, nil

}
