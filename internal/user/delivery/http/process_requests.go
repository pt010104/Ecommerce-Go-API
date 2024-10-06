package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	pkgErrors "github.com/pt010104/api-golang/pkg/errors"
	"github.com/pt010104/api-golang/pkg/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h handler) processSignupRequest(c *gin.Context) (signupReq, error) {
	ctx := c.Request.Context()

	var req signupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSignupRequest: invalid request")
		return signupReq{}, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSignupRequest: invalid request")
		return signupReq{}, err
	}

	return req, nil
}

func (h handler) processSignInRequest(c *gin.Context) (signinReq, error) {
	ctx := c.Request.Context()

	var req signinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSigninRequest: invalid request")
		return signinReq{}, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSignupRequest: invalid request")
		return signinReq{}, err
	}

	return req, nil
}

func (h handler) processDetailRequest(c *gin.Context) (string, models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.processDetailRequest: unauthorized")
		return "", models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		h.l.Errorf(ctx, "survey.delivery.http.handlers.processDetailRequest: invalid request")
		return "", models.Scope{}, errWrongQuery
	}

	sc := jwt.NewScope(payload)

	return id, sc, nil
}
