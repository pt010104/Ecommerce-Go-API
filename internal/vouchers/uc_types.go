package vouchers

import "time"

type CreateVoucherInput struct {
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
}
type ApplyVoucherInput struct {
	Code string
	//order amount is double type
	OrderAmount float64
}
