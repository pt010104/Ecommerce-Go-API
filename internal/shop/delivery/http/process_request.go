package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processRegisterRequest(c *gin.Context) (models.Scope, registerRequest, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDetailRequest: unauthorized")
		return models.Scope{}, registerRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processRegisterRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processRegisterRequest: invalid request")
		return models.Scope{}, req, err
	}

	sc := jwt.NewScope(payload)

	return sc, req, nil

}
