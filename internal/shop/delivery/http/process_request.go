package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateRequest(c *gin.Context) (models.Scope, registerRequest, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDetailRequest: unauthorized")
		return models.Scope{}, registerRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processCreateRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processCreateRequest: invalid request")
		return models.Scope{}, req, err
	}

	sc := jwt.NewScope(payload)

	return sc, req, nil

}

func (h handler) processGetRequest(c *gin.Context) (models.Scope, getShopRequest, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDetailRequest: unauthorized")
		return models.Scope{}, getShopRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req getShopRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetRequest: invalid request")
		return models.Scope{}, req, errWrongQuery
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetRequest: invalid request")
		return models.Scope{}, req, err
	}

	sc := jwt.NewScope(payload)

	return sc, req, nil

}

func (h handler) processDetailRequest(c *gin.Context) (models.Scope, string, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDetailRequest: unauthorized")
		return models.Scope{}, "", pkgErrors.NewUnauthorizedHTTPError()
	}

	id := c.Param("id")
	if id == "" {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processDetailRequest: invalid request")
		return models.Scope{}, "", errWrongBody
	}

	sc := jwt.NewScope(payload)

	return sc, id, nil
}
func (h handler) processDeleteShopRequest(c *gin.Context) (models.Scope, error) {
	ctx := c.Request.Context()
	payload, err := jwt.GetPayloadFromContext(ctx)
	if err != true {
		h.l.Errorf(ctx, " shop.Delivery.processDeleteShopRequest : ", err)
		return models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)
	return sc, nil
}
func (h handler) processUpdateShopRequest(c *gin.Context) (models.Scope, updateShopRequest, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processUpdateShopRequest: unauthorized")
		return models.Scope{}, updateShopRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req updateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processUpdateShopRequest: invalid request body")
		return models.Scope{}, req, errWrongBody
	}

	sc := jwt.NewScope(payload)

	return sc, req, nil
}
func (h handler) processGetShopIDByUserIDRequest(c *gin.Context) (models.Scope, GetShopIDByUserIDRequest, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetShopIDByUserIDRequest: unauthorized")
		return models.Scope{}, GetShopIDByUserIDRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req GetShopIDByUserIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetShopIDByUserIDRequestt: invalid request body")
		return models.Scope{}, req, errWrongBody
	}
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetShopIDByUserIDRequest: invalid request")
		return models.Scope{}, req, err
	}

	sc := jwt.NewScope(payload)

	return sc, req, nil
}
