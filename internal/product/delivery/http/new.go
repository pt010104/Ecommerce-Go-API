package http

import (
	"log"

	"github.com/pt010104/api-golang/internal/product"
)

type Handler interface {
}

type handler struct {
	l  log.Logger
	uc product.UseCase
}

func New(l log.Logger, uc product.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
