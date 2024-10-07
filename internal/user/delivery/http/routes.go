package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.POST("/sign-up", h.SignUp)
	r.POST("/sign-in", h.SignIn)

	r.Use(mw.Auth())
	r.GET("/:id", h.Detail)
	r.POST("/sign-out", h.SignOut)
}
