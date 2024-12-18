package http

import "github.com/gin-gonic/gin"

func MapPublicRoutes(r *gin.RouterGroup, h Handler) {
	r.GET("", h.GetProduct)
	r.GET("/:id", h.DetailProduct)
	r.GET("/get", h.GetAll)
}
