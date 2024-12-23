package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	MapCheckoutRouters(r, h, mw)
}

func MapCheckoutRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	group := r.Group("/checkout")
	group.Use(mw.Auth())
	group.POST("", h.CreateCheckout)
}
