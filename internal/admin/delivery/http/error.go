package http

import (
	"github.com/pt010104/api-golang/internal/admin"

	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errNoPermission = pkgErrors.NewHTTPError(130001, "you dont have permision to dothis")
	errWrongInput   = pkgErrors.NewHTTPError(130002, "category must have name and description")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case admin.ErrInvalidInput:
		return errWrongInput

	case admin.ErrNoPermission:
		return errNoPermission

	}

	return e
}
