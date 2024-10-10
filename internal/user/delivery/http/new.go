package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type Handler interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	Detail(c *gin.Context)
	SignOut(c *gin.Context)
	ForgetPasswordRequest(c *gin.Context)
	ResetPassword(c *gin.Context)
	VerifyRequest(c *gin.Context)
	VerifyUser(c *gin.Context)
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
