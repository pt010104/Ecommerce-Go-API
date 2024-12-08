package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	Upload(c *gin.Context)
}

type handler struct {
	l  log.Logger
	uc media.UseCase
}

func New(l log.Logger, uc media.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
