package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/checkout"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	Create(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc checkout.UseCase
}

func New(l log.Logger, uc checkout.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
