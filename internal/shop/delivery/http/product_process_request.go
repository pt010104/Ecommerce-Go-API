package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateProductRequest(c *gin.Context) (models.Scope, createProductReq, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "shop.http.delivery.hhtp.handler.processRequest: unauthorized")
		return models.Scope{}, createProductReq{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processCreateRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.validCreateRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil
}

func (h handler) processDetailProductRequest(c *gin.Context) (models.Scope, detailProductReq, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "shop.http.delivery.hhtp.handler.processRequest: unauthorized")
		return models.Scope{}, detailProductReq{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req detailProductReq
	if err := c.ShouldBindUri(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processCreateRequest: invalid request")
		return models.Scope{}, detailProductReq{}, errWrongBody
	}

	return sc, req, nil
}

func (h handler) processListProductRequest(c *gin.Context) (models.Scope, listProductRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDetailRequest: unauthorized")
		return models.Scope{}, listProductRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req listProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processListRequest: invalid request")
		return models.Scope{}, req, errWrongQuery
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil

}
func (h handler) processDeleteProductRequest(c *gin.Context) (models.Scope, deleteProductRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDeleteProductRequest: unauthorized")
		return models.Scope{}, deleteProductRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req deleteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processDeleteProductRequest: invalid request")
		return models.Scope{}, req, errWrongQuery
	}

	return sc, req, nil

}
func (h handler) processGetProductRequest(c *gin.Context) (models.Scope, getProductRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processGetRequest: unauthorized")
		return models.Scope{}, getProductRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req getProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetRequest: invalid request")
		return models.Scope{}, req, errWrongQuery
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil

}
