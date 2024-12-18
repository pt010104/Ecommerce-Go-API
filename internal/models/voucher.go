package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Voucher struct {
	ID                     primitive.ObjectID   `bson:"_id"`
	Name                   string               `bson:"name"`
	Description            string               `bson:"description,omitempty"`
	Code                   string               `bson:"code"`
	ValidFrom              time.Time            `bson:"valid_from"`
	ValidTo                time.Time            `bson:"valid_to"`
	UsageLimit             uint                 `bson:"usage_limit,omitempty"`
	ApplicableProductIDs   []primitive.ObjectID `bson:"applicable_product_ids,omitempty"`
	ApplicableCategorieIDs []primitive.ObjectID `bson:"applicable_category_ids,omitempty"`
	MinimumOrderAmount     float64              `bson:"minimum_order_amount,omitempty"`
	DiscountType           string               `bson:"discount_type"`
	DiscountAmount         float64              `bson:"discount_amount"`
	MaxDiscountAmount      float64              `bson:"max_discount_amount,omitempty"`
	UsedCount              int                  `bson:"used_count"`
	CreatedBy              primitive.ObjectID   `bson:"created_by"`
	Scope                  int                  `bson:"scope"` // 0: all, 1: shop
	ShopIDs                []primitive.ObjectID `bson:"shop_ids,omitempty"`
	CreatedAt              time.Time            `bson:"created_at"`
	UpdatedAt              time.Time            `bson:"updated_at"`
}

var (
	DiscountTypeFixed   = "fixed"
	DiscountTypePercent = "percent"

	ScopeAll  = 0
	ScopeShop = 1
)
