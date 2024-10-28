package http

import (
	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
)

type VerifyShopReq struct {
	ShopIDs []string `json:"ids" binding:"required"`
}

func (req VerifyShopReq) toInput() admin.VerifyShopInput {
	return admin.VerifyShopInput{
		ShopIDs: req.ShopIDs,
	}
}

type VerifyShopResp struct {
	ShopID     string `json:"shop_id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

func (h handler) toResList(shops []models.Shop) []VerifyShopResp {
	var resList []VerifyShopResp
	for _, shop := range shops {
		res := VerifyShopResp{
			ShopID:     shop.ID.Hex(),
			Name:       shop.Name,
			IsVerified: shop.IsVerified,
		}
		resList = append(resList, res)
	}
	return resList
}
