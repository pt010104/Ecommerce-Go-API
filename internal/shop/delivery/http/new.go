package http

import (
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
}

type handler struct {
	l  log.Logger
	uc shop.UseCase
}

func New(l log.Logger, uc shop.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
