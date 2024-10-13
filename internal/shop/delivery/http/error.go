package http

import (
	"github.com/pt010104/api-golang/internal/shop"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongPaginationQuery = pkgErrors.NewHTTPError(130001, "Wrong pagination query")
	errWrongQuery           = pkgErrors.NewHTTPError(130002, "Wrong query")
	errWrongBody            = pkgErrors.NewHTTPError(130003, "Wrong body")
	errWrongHeader          = pkgErrors.NewHTTPError(130004, "Wrong header")

	ErrInvalidPhone = pkgErrors.NewHTTPError(130005, "Invalid phone")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case shop.ErrInvalidPhone:
		return ErrInvalidPhone
	}

	return e
}
