package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	SignUp(c *gin.Context)
}
type handler struct {
	l  log.Logger
	uc user.UseCase
}

func New(l log.Logger, uc user.UseCase) Handler {
	return &handler{
		l:  l,
		uc: uc,
	}
}
