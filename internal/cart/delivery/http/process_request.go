package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
)

func (h handler) processUpdateCartRequest(c *gin.Context) (models.Scope, UpdateCartRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processUpdateCartRequest: unauthorized")
		return models.Scope{}, UpdateCartRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processUpdateCartRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCartRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil

}

func (h handler) processaddToCartRequest(c *gin.Context) (models.Scope, addToCartRequest, error) {
	ctx := c.Request.Context()

	sc, ok := jwt.GetScopeFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processaddToCartRequest: unauthorized")
		return models.Scope{}, addToCartRequest{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	var req addToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processaddToCartRequest: invalid request")
		return models.Scope{}, req, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.handler.processaddToCartRequest: invalid request")
		return models.Scope{}, req, err
	}

	return sc, req, nil
}

// func (h handler) processListCartRequest(c *gin.Context) (models.Scope, ListCartRequest, error) {
// 	ctx := c.Request.Context()

// 	sc, ok := jwt.DetailScopeFromContext(ctx)
// 	if !ok {
// 		h.l.Errorf(ctx, "survey.delivery.http.handler.processListCartRequest: unauthorized")
// 		return models.Scope{}, ListCartRequest{}, pkgErrors.NewUnauthorizedHTTPError()
// 	}

// 	var req ListCartRequest
// 	req.UserID = c.DetailHeader("x-client-id")
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		h.l.Errorf(ctx, "cart.delivery.http.handler.processListCartRequest: invalid request")
// 		return models.Scope{}, req, errWrongBody
// 	}

// 	if err := req.validate(); err != nil {
// 		h.l.Errorf(ctx, "cart.delivery.http.handler.processCreateCartRequest: invalid request")
// 		return models.Scope{}, req, err
// 	}

// 	return sc, req, nil

// }
