package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateCheckoutRequest(c *gin.Context) (models.Scope, CreateCheckoutRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCheckoutRequest: unauthorized")
		return models.Scope{}, CreateCheckoutRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCheckoutRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCheckoutRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil

}

func (h *handler) processCreateOrderRequest(c *gin.Context) (models.Scope, CreateOrderRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCheckoutRequest: unauthorized")
		return models.Scope{}, CreateOrderRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processCreateOrderRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processCreateOrderRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil
}

func (h *handler) processListOrderRequest(c *gin.Context) (models.Scope, ListOrderRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "order.delivery.http.handler.processListOrderRequest: unauthorized")
		return models.Scope{}, ListOrderRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req ListOrderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processListOrderRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processListOrderRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil
}
