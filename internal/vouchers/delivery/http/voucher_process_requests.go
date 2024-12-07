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
