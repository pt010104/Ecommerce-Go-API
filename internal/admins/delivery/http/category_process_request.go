package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateCategoryRequest(c *gin.Context) (createCategoryReq, models.Scope, error) {
	ctx := c.Request.Context()
	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "admins.http.delivery.hhtp.handler.processRequest: unauthorized")
		return createCategoryReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return createCategoryReq{}, models.Scope{}, errWrongInput
	}
	sc := jwt.NewScope(payload)
	role, exist := c.Get("role")
	if !exist {
		h.l.Errorf(ctx, "admins.http.delivery.handler.processrequest.getrole:")
		return createCategoryReq{}, models.Scope{}, errWrongInput
	}
	roleInt, ok := role.(int)
	if !ok {
		h.l.Errorf(ctx, "admins.http.delivery.handler.processRequest.getRole: role is not an int")
		return createCategoryReq{}, models.Scope{}, errWrongInput
	}

	sc.Role = roleInt
	return req, sc, nil

}
