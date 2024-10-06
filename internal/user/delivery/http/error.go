package http

import (
	"github.com/pt010104/api-golang/internal/user"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongPaginationQuery = pkgErrors.NewHTTPError(120001, "Wrong pagination query")
	errWrongQuery           = pkgErrors.NewHTTPError(120002, "Wrong query")
	errWrongBody            = pkgErrors.NewHTTPError(120003, "Wrong body")

	ErrEmailExisted = pkgErrors.NewHTTPError(120004, "email has already been registered")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case user.ErrEmailExisted:
		return ErrEmailExisted
	}

	return e
}
