package cart

import "errors"

var (
	ErrNotEnoughQuantity = errors.New("not enough quantity")
	ErrInvalidProductID  = errors.New("invalid product ID")
	ErrCartNotFound      = errors.New("cart not found")
	ErrUserMismatch      = errors.New("user ID does not match cart's user ID")
	ErrEmptyItemList     = errors.New("new item list cannot be empty")
	ErrInvalidQuantity   = errors.New("invalid quantity for product")
	ErrWrongBody         = errors.New("Wrong body")
)
