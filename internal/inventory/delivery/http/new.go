package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/inventory"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	Create(c *gin.Context)
}
type handler struct {
	l  log.Logger
	uc inventory.UseCase
}

func New(l log.Logger, uc inventory.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
