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

	return req, nil
}
