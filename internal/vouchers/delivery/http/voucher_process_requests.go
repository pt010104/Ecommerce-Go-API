package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateVoucherRequest(c *gin.Context) (createVoucherReq, models.Scope, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "admin.http.delivery.hhtp.handler.processRequest: unauthorized")
		return createVoucherReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createVoucherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return createVoucherReq{}, models.Scope{}, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request %v", err)
		return createVoucherReq{}, models.Scope{}, errWrongBody
	}

	return req, sc, nil

}

func (h handler) processDetailVoucherRequest(c *gin.Context) (DetailVoucherReq, models.Scope, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "voucher.http.delivery.hhtp.handler.processRequest: unauthorized")
		return DetailVoucherReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req DetailVoucherReq
	if err := c.ShouldBindUri(&req); err != nil {
		h.l.Errorf(ctx, "voucher.delivery.http.handler.processDetailRequest: invalid request")
		return DetailVoucherReq{}, models.Scope{}, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "voucher.delivery.http.handler.processDetailRequest: invalid request %v", err)
		return DetailVoucherReq{}, models.Scope{}, errWrongBody
	}

	return req, sc, nil

}
func (h handler) processListVoucherRequest(c *gin.Context) (ListVoucherReq, models.Scope, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "voucher.http.delivery.hhtp.handler.processRequest: unauthorized")
		return ListVoucherReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req ListVoucherReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "voucher.delivery.http.handler.processListRequest: invalid request")
		return ListVoucherReq{}, models.Scope{}, errWrongBody
	}
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "voucher.delivery.http.handler.processListRequest: invalid request %v", err)
		return ListVoucherReq{}, models.Scope{}, errWrongBody
	}

	return req, sc, nil

}

func (h handler) processApplyVoucherRequest(c *gin.Context) (applyVoucherReq, models.Scope, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "voucher.http.delivery.http.handler.processApplyVoucherRequest: unauthorized")
		return applyVoucherReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req applyVoucherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "voucher.delivery.http.handler.processApplyVoucherRequest: invalid request")
		return applyVoucherReq{}, models.Scope{}, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "voucher.delivery.http.handler.processApplyVoucherRequest: invalid request %v", err)
		return applyVoucherReq{}, models.Scope{}, errWrongBody
	}

	return req, sc, nil
}
