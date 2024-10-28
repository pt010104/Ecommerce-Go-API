package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processVerifyShopRequest(c *gin.Context) (VerifyShopReq, models.Scope, error) {

	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if ok != true {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: cannot get scope from context")
		return VerifyShopReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req VerifyShopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return VerifyShopReq{}, models.Scope{}, errWrongInput
	}

	return req, sc, nil
}
