package http

import "github.com/gin-gonic/gin"

func MapRouters(r *gin.RouterGroup, h Handler) {
	r.POST("/signUp", h.SignUp)
}
