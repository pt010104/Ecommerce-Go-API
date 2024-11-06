package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h handler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Create: %v", err)
		response.Error(c, err)
		return
	}

	product, _, err := h.uc.CreateProduct(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Create: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, product)
}
