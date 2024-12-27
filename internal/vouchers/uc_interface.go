package vouchers

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateVoucher(ctx context.Context, sc models.Scope, input CreateVoucherInput) (models.Voucher, error)
	Detail(ctx context.Context, sc models.Scope, input DetailVoucherInput) (models.Voucher, error)
	List(ctx context.Context, sc models.Scope, opt GetVoucherFilter) ([]models.Voucher, error)
	ApplyVoucher(ctx context.Context, sc models.Scope, input ApplyVoucherInput) (models.Voucher, float64, float64, error)
}
