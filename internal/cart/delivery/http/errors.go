package http

import (
	"github.com/pt010104/api-golang/internal/cart"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongPaginationQuery  = pkgErrors.NewHTTPError(120001, "Wrong pagination query")
	errWrongQuery            = pkgErrors.NewHTTPError(120002, "Wrong query")
	errWrongBody             = pkgErrors.NewHTTPError(120003, "Wrong body")
	errWrongHeader           = pkgErrors.NewHTTPError(120004, "Wrong header")
	ErrMismatchedUser        = pkgErrors.NewHTTPError(120004, "This cart belong to another user")
	ErrEmailExisted          = pkgErrors.NewHTTPError(120005, "email has already been registered")
	ErrUserNotVerified       = pkgErrors.NewHTTPError(120006, "user not verified")
	ErrInvalidPasswordFormat = pkgErrors.NewHTTPError(120007, "password must contain number and digits")
	ErrInvalidEmailFormat    = pkgErrors.NewHTTPError(120007, "invalid email")
	ErrInvalidNameFormat     = pkgErrors.NewHTTPError(120007, "invalid name")
	ErrStockNotEnough        = pkgErrors.NewHTTPError(120008, "not enough stock")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case cart.ErrWrongBody:
		return errWrongBody
	case cart.ErrUserMismatch:
		return ErrMismatchedUser
	case cart.ErrNotEnoughStock:
		return ErrStockNotEnough
	}

	return e
}
