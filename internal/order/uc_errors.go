package order

import "errors"

var (
	ErrCartNotFound          = errors.New("cart not found")
	ErrProductNotFoundInCart = errors.New("product not found in cart")
	ErrProductNotFound       = errors.New("product not found")
	ErrProductNotEnoughStock = errors.New("product not enough stock")
	ErrCheckoutStatusInvalid = errors.New("checkout status invalid")
	ErrCheckoutExpired       = errors.New("checkout expired")
	ErrTooManyRetries        = errors.New("too many retries due to concurrent modifications")
	ErrRedisNotFound         = errors.New("redis: nil")
)
