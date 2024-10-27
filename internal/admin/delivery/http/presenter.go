package http

import (
	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
)

type VerifyShopReq struct {
	ShopID []string `json:"shop_id" binding:"required"`
}

func (req VerifyShopReq) toInput() admin.VerifyShopInput {
	return admin.VerifyShopInput{
		ShopID: req.ShopID,
	}
}

type VerifyShopResp struct {
	ShopID     string `json:"shop_id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}

func (resp VerifyShopResp) toResList(shops []models.Shop) []VerifyShopResp {
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
