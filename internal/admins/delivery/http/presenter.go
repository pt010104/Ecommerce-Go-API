package http

import (
	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
)

type VerifyShopReq struct {
	ShopIDs []string `json:"ids" binding:"required"`
}

func (req VerifyShopReq) toInput() admins.VerifyShopInput {
	return admins.VerifyShopInput{
		ShopIDs: req.ShopIDs,
	}
}

type verifyShopResp struct {
	ShopID     string `json:"shop_id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

func (h handler) toResList(shops []models.Shop) []verifyShopResp {
	var resList []verifyShopResp
	for _, shop := range shops {
		res := verifyShopResp{
			ShopID:     shop.ID.Hex(),
			Name:       shop.Name,
			IsVerified: shop.IsVerified,
		}
		resList = append(resList, res)
	}
	return resList
}
