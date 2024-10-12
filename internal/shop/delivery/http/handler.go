package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processRegisterRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.Register: %v", err)
		response.Error(c, err)
		return
	}

	shop, err := h.uc.Register(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.Register: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, shop)
}
