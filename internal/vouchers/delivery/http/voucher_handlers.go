package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

func (h handler) CreateVoucher(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateVoucherRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processCreateRequest: %v", err)
		response.Error(c, err)
		return
	}
	u, err := h.uc.CreateVoucher(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, u)
}
func (h handler) DetailVoucher(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processDetailVoucherRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processDetailRequest: %v", err)
		response.Error(c, err)
		return
	}
	u, err := h.uc.Detail(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newDetailResponse(u))
}
func (h handler) ListVoucher(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processListVoucherRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processListRequest: %v", err)
		response.Error(c, err)
		return
	}
	u, err := h.uc.List(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newListResponse(u))
}
