package http

import (
	"github.com/pt010104/api-golang/internal/user"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongPaginationQuery = pkgErrors.NewHTTPError(120001, "Wrong pagination query")
	errWrongQuery           = pkgErrors.NewHTTPError(120002, "Wrong query")
	errWrongBody            = pkgErrors.NewHTTPError(120003, "Wrong body")
	errWrongHeader          = pkgErrors.NewHTTPError(120004, "Wrong header")

	ErrEmailExisted    = pkgErrors.NewHTTPError(120005, "email has already been registered")
	ErrUserNotVerified = pkgErrors.NewHTTPError(120006, "user not verified")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case user.ErrEmailExisted:
		return ErrEmailExisted
	case user.ErrUserNotVerified:
		return ErrUserNotVerified
	}

	return e
}
