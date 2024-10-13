package shop

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")

	ErrInvalidPhone = errors.New("invalid phone")
	ErrShopExist    = errors.New("shop exist")
)
