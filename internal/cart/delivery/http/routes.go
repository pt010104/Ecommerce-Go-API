package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.POST("/create-cart", h.Create)
	r.POST("/update-cart", h.Update)
	r.GET("", h.List)
	r.GET("get-cart", h.Get)

}
