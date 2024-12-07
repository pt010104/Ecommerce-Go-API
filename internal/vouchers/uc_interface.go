package vouchers

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateVoucher(ctx context.Context, sc models.Scope, input CreateVoucherInput) (models.Voucher, error)
}
