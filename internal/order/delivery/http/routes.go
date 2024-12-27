package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	MapOrderRouters(r, h, mw)
	MapCheckoutRouters(r, h, mw)
}

func MapOrderRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.Use(mw.Auth())
	r.POST("", h.CreateOrder)
	r.GET("", h.ListOrder)
}

func MapCheckoutRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	group := r.Group("/checkout")
	group.Use(mw.Auth())
	group.POST("", h.CreateCheckout)
}
