package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	CreateCategory(c *gin.Context)
	VerifyShop(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc admin.UseCase
}

func New(l log.Logger, uc admin.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
