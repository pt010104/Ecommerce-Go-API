package http

import (
	"github.com/pt010104/api-golang/internal/vouchers"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errNoPermission = pkgErrors.NewHTTPError(140001, "you dont have permision to dothis")
	errWrongBody    = pkgErrors.NewHTTPError(140002, "Wrong body")

	errExistCode                 = pkgErrors.NewHTTPError(140003, "code already exist")
	errShopNotFound              = pkgErrors.NewHTTPError(140004, "shop not found")
	ErrVoucherNotFound           = pkgErrors.NewHTTPError(140005, "voucher not found")
	ErrVoucherExpired            = pkgErrors.NewHTTPError(140006, "voucher expired")
	ErrVoucherMinimumOrderAmount = pkgErrors.NewHTTPError(140007, "voucher minimum order amount")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case vouchers.ErrExistCode:
		return errExistCode

	case vouchers.ErrShopNotFound:
		return errShopNotFound
	case vouchers.ErrVoucherNotFound:
		return ErrVoucherNotFound
	case vouchers.ErrVoucherExpired:
		return ErrVoucherExpired
	case vouchers.ErrVoucherMinimumOrderAmount:
		return ErrVoucherMinimumOrderAmount

	}

	return e
}
