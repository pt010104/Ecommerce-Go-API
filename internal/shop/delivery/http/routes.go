package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())

	r.POST("", h.Create)
	r.GET("", h.Get)
	r.POST("/delete", h.Delete)

}
