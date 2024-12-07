package vouchers

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrExistCode    = errors.New("code already exist")

	ErrShopNotFound = errors.New("shop not found")
)
