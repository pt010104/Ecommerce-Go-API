package http

import "github.com/gin-gonic/gin"

func MapRouters(r *gin.RouterGroup, h Handler) {
	r.POST("/sign_up", h.SignUp)
	r.POST("/sign_in", h.SignIn)
}
