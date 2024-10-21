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

	return req, nil
}
func (h handler) processForgetPasswordRequest(c *gin.Context) (forgetPasswordReq, error) {
	ctx := c.Request.Context()

	var req forgetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSigninRequest: invalid request")
		return forgetPasswordReq{}, errWrongBody
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

	req.SessionID = c.GetHeader("session-id")

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
func (h handler) processLogOutRequest(c *gin.Context) (models.Scope, error) {
	ctx := c.Request.Context()

	payload, ok := jwt.GetPayloadFromContext(ctx)
	if !ok {
		h.l.Errorf(ctx, "survey.delivery.http.handler.procesLogOutRequest: unauthorized")
		return models.Scope{}, pkgErrors.NewUnauthorizedHTTPError()
	}

	sc := jwt.NewScope(payload)
	return sc, nil

}

func (h handler) processResetPasswordRequest(c *gin.Context) (resetPasswordReq, error) {
	ctx := c.Request.Context()

	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processResetPasswordRequest: invalid request")
		return resetPasswordReq{}, errWrongBody
	}

	req.Token = c.Query("token")
	userID, exist := c.Get("userID")
	if !exist {
		return resetPasswordReq{}, errWrongBody
	}
	req.UserID = userID.(string)
	return req, nil
}
func (h handler) processVerifyEmailRequest(c *gin.Context) (verifyRequestReq, error) {
	ctx := c.Request.Context()

	var req verifyRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processVerifyEmailRequest: invalid request")
		return verifyRequestReq{}, errWrongBody
	}

	return req, nil

}
func (h handler) processVerifyUserRequest(c *gin.Context) (verifyUserReq, error) {
	ctx := c.Request.Context()

	var req verifyUserReq

	req.Token = c.Query("token")
	userID, exist := c.Get("userID")
	if !exist {
		return verifyUserReq{}, errWrongBody
	}
	req.UserID = userID.(string)
	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processVerifyUserRequest: invalid request")
		return verifyUserReq{}, err
	}

	return req, nil

}
func (h handler) processDistributeNewTokenRequest(c *gin.Context) (distributeNewTokenReq, error) {
	ctx := c.Request.Context()
	var req distributeNewTokenReq

	req.SessionID = c.GetHeader("session-id")
	req.UserId = c.GetHeader("x-client-id")
	req.RefreshToken = c.GetHeader("refresh-token")

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processDistributeNewTokenRequest: invalid request")
		return distributeNewTokenReq{}, errWrongHeader
	}

	return req, nil
}
