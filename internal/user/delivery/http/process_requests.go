package http

import (
	"github.com/gin-gonic/gin"
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

	exists, err := h.uc.EmailExisted(ctx, req.Email)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSignupRequest: error checking email existence")
		return signupReq{}, err
	}

	if exists {
		h.l.Errorf(ctx, "user.delivery.http.handler.processSignupRequest: email existed")
		return signupReq{}, errEmailExisted
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
