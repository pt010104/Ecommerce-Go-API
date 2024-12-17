package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	CreateVoucher(c *gin.Context)
	DetailVoucher(c *gin.Context)
	ListVoucher(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc vouchers.UseCase
}

func New(l log.Logger, uc vouchers.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
