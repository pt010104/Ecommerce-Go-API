package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h handler) VerifyShop(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processVerifyShopRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "admin.delivery.http.handler.Create.processCreateRequest: %v", err)
		response.Error(c, err)
		return
	}
	h.l.Debugf(ctx, "role:", sc.Role)
	s, err := h.uc.VerifyShop(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "admin.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}
	h.l.Debugf(ctx, "role:", sc.Role)
	response.OK(c, s)
}
