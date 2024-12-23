package http

import (
	"github.com/pt010104/api-golang/internal/order"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongBody             = pkgErrors.NewHTTPError(120003, "Wrong body")
	errStockNotEnough        = pkgErrors.NewHTTPError(120004, "Stock not enough")
	errProductNotFound       = pkgErrors.NewHTTPError(120005, "Product not found")
	errProductNotFoundInCart = pkgErrors.NewHTTPError(120006, "Product not found in cart")
	errCartNotFound          = pkgErrors.NewHTTPError(120007, "Cart not found")
	errTooManyRetries        = pkgErrors.NewHTTPError(120008, "Too many retries due to concurrent modifications")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case order.ErrProductNotEnoughStock:
		return errStockNotEnough
	case order.ErrRedisNotFound:
		return errStockNotEnough
	case order.ErrProductNotFound:
		return errProductNotFound
	case order.ErrProductNotFoundInCart:
		return errProductNotFoundInCart
	case order.ErrCartNotFound:
		return errCartNotFound
	case order.ErrTooManyRetries:
		return errTooManyRetries
	}

	return e
}
