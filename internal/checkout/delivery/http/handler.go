package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h *handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "checkout.delivery.http.handler.Create: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.uc.Create(ctx, sc, req.ProductIDs)
	if err != nil {
		h.l.Errorf(ctx, "checkout.delivery.http.handler.Create: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newCreateCheckoutResponse(o))
}
