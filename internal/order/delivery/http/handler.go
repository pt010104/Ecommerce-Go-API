package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h *handler) CreateCheckout(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateCheckoutRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.CreateCheckout: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.uc.CreateCheckout(ctx, sc, req.ProductIDs)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.CreateCheckout: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newCreateCheckoutCheckoutResponse(o))
}
