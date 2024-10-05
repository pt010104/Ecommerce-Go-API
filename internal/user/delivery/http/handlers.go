package http

import (
	"github.com/gin-gonic/gin"

	"github.com/pt010104/api-golang/pkg/response"
)

func (h handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := h.processSignupRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signup.processSignupRequest: %v", err)
		response.Error(c, err)
		return
	}

	u, err := h.uc.CreateUser(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signup.uc.CreateUser: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, u)
}

func (h handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processSignInRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signup.processSigninRequest: %v", err)
		response.Error(c, err)
		return
	}

	token, err := h.uc.SignIn(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signin.uc.signIn: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, token)

}
