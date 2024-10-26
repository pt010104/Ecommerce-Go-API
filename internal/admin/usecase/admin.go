package usecase

import (
	"context"
	"fmt"
	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

func (uc implUsecase) VerifyShop(ctx context.Context, sc models.Scope, input admin.VerifyShopInput) ([]models.Shop, error) {
	if sc.Role == 0 {
		return []models.Shop{}, admin.ErrNoPermission
	}
	var listShop []models.Shop
	for i := 0; i < len(input.ShopID); i++ {
		fmt.Println("Updating ShopID:", input.ShopID[i])
		shop, err := uc.shopUc.Update(ctx, models.Scope{}, shop.UpdateInput{
			ShopID:     input.ShopID[i],
			IsVerified: true,
		})
		if err != nil {
			uc.l.Errorf(ctx, "admin.usecase.Verifyshop.shopUC.update:", err)
			return []models.Shop{}, err
		}
		listShop = append(listShop, shop)
	}

	return listShop, nil

}
