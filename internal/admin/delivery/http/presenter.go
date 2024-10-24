package http

import "github.com/pt010104/api-golang/internal/admin"

type VerifyShopReq struct {
	ShopID string `json:"shop_id" binding:"required"`
}

func (req VerifyShopReq) toInput() admin.VerifyShopInput {
	return admin.VerifyShopInput{
		ShopID: req.ShopID,
	}
}
