package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processCreateCategoryRequest(c *gin.Context) (createCategoryReq, models.Scope, error) {
	ctx := c.Request.Context()
	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "admins.http.delivery.hhtp.handler.processRequest: unauthorized")
		return createCategoryReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return createCategoryReq{}, models.Scope{}, errWrongInput
	}

	return req, sc, nil

}

func (h handler) processListCategoryRequest(c *gin.Context) (listCatagoryReq, models.Scope, error) {
	ctx := c.Request.Context()
	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "admins.http.delivery.hhtp.handler.processRequest: unauthorized")
		return listCatagoryReq{}, models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req listCatagoryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return listCatagoryReq{}, models.Scope{}, errWrongBody
	}

	return req, sc, nil

}
