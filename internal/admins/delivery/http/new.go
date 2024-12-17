package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	CreateCategory(c *gin.Context)
	VerifyShop(c *gin.Context)
	ListCategory(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc admins.UseCase
}

var _ Handler = &handler{}

func New(l log.Logger, uc admins.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
