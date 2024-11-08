package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
	CreateProduct(c *gin.Context)
	DetailProduct(c *gin.Context)
	ListProduct(c *gin.Context)
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
