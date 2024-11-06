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
		h.l.Errorf(ctx, "admin.http.delivery.hhtp.handler.processRequest: unauthorized")
		return models.Scope{}, createProductReq{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.handler.processCreateRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	return sc, req, nil
}
