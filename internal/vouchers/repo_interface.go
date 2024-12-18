package vouchers

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repository interface {
	CreateVoucher(ctx context.Context, sc models.Scope, opt CreateVoucherOption) (models.Voucher, error)
	DetailVoucher(ctx context.Context, sc models.Scope, opt DetailVoucherOption) (models.Voucher, error)
	ListVoucher(ctx context.Context, sc models.Scope, opt GetVoucherFilter) ([]models.Voucher, error)
	UpdateVoucher(ctx context.Context, sc models.Scope, option UpdateVoucherOption) (models.Voucher, error)
}
