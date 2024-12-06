package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateCartRequest(c *gin.Context) (models.Scope, CreateCartRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processCreateCartRequest: unauthorized")
		return models.Scope{}, CreateCartRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req CreateCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCartRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCartRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil

}
