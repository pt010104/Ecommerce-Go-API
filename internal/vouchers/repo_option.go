package vouchers

import (
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type VoucherData struct {
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
	UsedCount              int
}

type CreateVoucherOption struct {
	Data VoucherData
}
type GetVoucherFilter struct {
	ValidFrom              *time.Time
	ValidTo                *time.Time
	Scope                  int
	Codes                  []string
	IDs                    []string
	ShopIDs                []string
	ApplicableCategorieIDs []string
	ApplicableProductIDs   []string
}

type GetVoucherOption struct {
	Filter GetVoucherFilter
	Pag    paginator.Paginator
}

type UpdateVoucherOption struct {
	Model models.Voucher
	Data  VoucherData
}

type DetailVoucherOption struct {
	ID   string
	Code string
}
