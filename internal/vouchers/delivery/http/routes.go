package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {

	r.Use(mw.Auth())
	r.GET("by-code/:code", h.DetailVoucher)
	r.GET("by-id/:id", h.DetailVoucher)
	r.GET("", h.ListVoucher)
	r.POST("apply", h.ApplyVoucher)

	r.Use(mw.Auth(), mw.AuthShop())
	r.POST("", h.CreateVoucher)

}
