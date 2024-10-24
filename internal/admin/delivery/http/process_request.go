package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processVerifyShopRequest(c *gin.Context) (VerifyShopReq, models.Scope, error) {

	ctx := c.Request.Context()
	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "admin.http.delivery.hhtp.handler.processRequest: unauthorized")
		return VerifyShopReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}
	var req VerifyShopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return VerifyShopReq{}, models.Scope{}, errWrongInput
	}
	sc := jwt.NewScope(payload)
	role, exist := c.Get("role")
	if !exist {
		h.l.Errorf(ctx, "admin.http.delivery.handler.processrequest.getrole:")
		return VerifyShopReq{}, models.Scope{}, errWrongInput
	}
	roleInt, ok := role.(int)
	if !ok {
		h.l.Errorf(ctx, "admin.http.delivery.handler.processRequest.getRole: role is not an int")
		return VerifyShopReq{}, models.Scope{}, errWrongInput
	}

	sc.Role = roleInt
	return req, sc, nil
}
