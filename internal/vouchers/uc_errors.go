package vouchers

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrExistCode    = errors.New("code already exist")

	ErrShopNotFound              = errors.New("shop not found")
	ErrVoucherNotFound           = errors.New("voucher not found")
	ErrVoucherExpired            = errors.New("voucher expired")
	ErrVoucherMinimumOrderAmount = errors.New("voucher minimum order amount")
)
