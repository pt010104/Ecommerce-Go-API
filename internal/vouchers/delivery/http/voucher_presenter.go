package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
)

type createVoucherReq struct {
	Name                   string   `json:"name" binding:"required"`
	Code                   string   `json:"code" binding:"required"`
	ValidFrom              string   `json:"valid_from" binding:"required"`
	ValidTo                string   `json:"valid_to" binding:"required"`
	DiscountType           string   `json:"discount_type" binding:"required"`
	DiscountAmount         float64  `json:"discount_amount" binding:"required"`
	Description            string   `json:"description"`
	UsageLimit             uint     `json:"usage_limit"`
	ApplicableProductIDs   []string `json:"applicable_product_ids"`
	ApplicableCategorieIDs []string `json:"applicable_category_ids"`
	MinimumOrderAmount     float64  `json:"minimum_order_amount"`
	MaxDiscountAmount      float64  `json:"max_discount_amount"`
	ShopIDs                []string `json:"shop_ids"`
}

func (req createVoucherReq) validate() error {
	if _, err := util.StrToDateTime(req.ValidFrom); err != nil {
		return errWrongBody
	}
	if _, err := util.StrToDateTime(req.ValidTo); err != nil {
		return errWrongBody
	}

	if req.DiscountType != models.DiscountTypePercent && req.DiscountType != models.DiscountTypeFixed {
		return errWrongBody
	}

	if len(req.ApplicableProductIDs) > 0 {
		for _, id := range req.ApplicableProductIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	if len(req.ApplicableCategorieIDs) > 0 {
		for _, id := range req.ApplicableCategorieIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	if len(req.ShopIDs) > 0 {
		for _, id := range req.ShopIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	return nil
}

func (req createVoucherReq) toInput() vouchers.CreateVoucherInput {
	validFrom, _ := util.StrToDateTime(req.ValidFrom)
	validTo, _ := util.StrToDateTime(req.ValidTo)

	return vouchers.CreateVoucherInput{
		Name:                   req.Name,
		Description:            req.Description,
		Code:                   req.Code,
		ValidFrom:              validFrom,
		ValidTo:                validTo,
		UsageLimit:             req.UsageLimit,
		ApplicableProductIDs:   req.ApplicableProductIDs,
		ApplicableCategorieIDs: req.ApplicableCategorieIDs,
		MinimumOrderAmount:     req.MinimumOrderAmount,
		DiscountType:           req.DiscountType,
		DiscountAmount:         req.DiscountAmount,
		MaxDiscountAmount:      req.MaxDiscountAmount,
		ShopIDs:                req.ShopIDs,
	}
}

type detailVoucherResp struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Code                   string   `json:"code"`
	ValidFrom              string   `json:"valid_from"`
	ValidTo                string   `json:"valid_to"`
	DiscountType           string   `json:"discount_type"`
	DiscountAmount         float64  `json:"discount_amount"`
	Description            string   `json:"description,omitempty"`
	UsageLimit             uint     `json:"usage_limit"`
	ApplicableProductIDs   []string `json:"applicable_product_ids,omitempty"`
	ApplicableCategorieIDs []string `json:"applicable_category_ids,omitempty"`
	MinimumOrderAmount     float64  `json:"minimum_order_amount"`
	MaxDiscountAmount      float64  `json:"max_discount_amount"`
	ShopIDs                []string `json:"shop_ids"`
	CreatedAt              string   `json:"created_at"`
}
