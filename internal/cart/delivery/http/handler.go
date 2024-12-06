package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Create: %v", err)
		response.Error(c, err)
		return
	}

	cart, err := h.uc.Create(sc, ctx, req.toInput(), req.toItemInput())
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Create: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, cart)
}
