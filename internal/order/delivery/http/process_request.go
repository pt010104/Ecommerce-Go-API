package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/mongo"
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

func (h *handler) processListOrderShopRequest(c *gin.Context) (models.Scope, ListOrderShopRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "order.delivery.http.handler.processListOrderShopRequest: unauthorized")
		return models.Scope{}, ListOrderShopRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req ListOrderShopRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processListOrderShopRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processListOrderShopRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil
}

func (h *handler) processUpdateOrderRequest(c *gin.Context) (models.Scope, string, UpdateOrderRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "order.delivery.http.handler.processUpdateOrderRequest: unauthorized")
		return models.Scope{}, "", UpdateOrderRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	orderID := c.Param("order_id")
	if !mongo.IsObjectID(orderID) {
		h.l.Errorf(ctx, "order.delivery.http.handler.processUpdateOrderRequest: invalid order_id")
		return models.Scope{}, "", UpdateOrderRequest{}, errWrongBody
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processUpdateOrderRequest: invalid request")
		return models.Scope{}, "", req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.processUpdateOrderRequest: invalid request")
		return models.Scope{}, "", req, err
	}

	return sc, orderID, req, nil
}
