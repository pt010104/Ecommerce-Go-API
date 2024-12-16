package cart

import "errors"

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrUserMismatch    = errors.New("user ID does not match cart's user ID")
	ErrEmptyItemList   = errors.New("new item list cannot be empty")
	ErrInvalidQuantity = errors.New("invalid quantity for product")
	ErrWrongBody       = errors.New("Wrong body")
	ErrInvalidCartItem = errors.New("invalid cart item")
	ErrNotEnoughStock  = errors.New("not enough stock")
	ErrRequiredField   = errors.New("required field")
)
