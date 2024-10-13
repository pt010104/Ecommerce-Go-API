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
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processGetRequest: invalid request")
		return models.Scope{}, req, errWrongBody
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
