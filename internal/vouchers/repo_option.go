package vouchers

import "time"

type CreateVoucherOption struct {
	Name                   string
	ShopIDs                []string
	Description            string
	Code                   string
	ValidFrom              time.Time
	ValidTo                time.Time
	UsageLimit             uint
	ApplicableProductIDs   []string
	ApplicableCategorieIDs []string
	MinimumOrderAmount     float64
	DiscountType           string
	DiscountAmount         float64
	MaxDiscountAmount      float64
	Scope                  int
	CreatedBy              string
}
