package http

import (
	"github.com/pt010104/api-golang/internal/media"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
)

var (
	errWrongBody = pkgErrors.NewHTTPError(150001, "Wrong body")
	errNoFiles   = pkgErrors.NewHTTPError(150002, "No files provided")
)

func (h handler) mapErrors(e error) error {
	switch e {
	case media.ErrRequireField:
		return errWrongBody
	}

	return e
}
