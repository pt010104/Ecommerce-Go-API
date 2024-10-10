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
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, u)
}

func (h handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processSignInRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signin.processSigninRequest: %v", err)
		response.Error(c, err)
		return
	}

	output, err := h.uc.SignIn(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signin.uc.signIn: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	signInResp := h.newSignInResp(output)
	response.OK(c, signInResp)

}

func (h handler) ForgetPasswordRequest(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processForgetPasswordRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.Signup.processForgetPassRequest: %v", err)
		response.Error(c, err)
		return
	}

	h.uc.ForgetPasswordRequest(ctx, sc.Email)

	response.OK(c, nil)
}

func (h handler) SignOut(c *gin.Context) {
	ctx := c.Request.Context()
	sc, err := h.processLogOutRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.Signup.processSignOutRequest: %v", err)
		response.Error(c, err)
		return
	}
	h.uc.LogOut(ctx, sc)

}
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.DetailUser.processDetailUserRequest: %v", err)
		response.Error(c, err)
		return
	}

	u, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.DetailUser.uc.DetailUser: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newDetailResp(u))
}

func (h handler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processResetPasswordRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword.processResetPasswordRequest: %v", err)
		response.Error(c, err)
		return
	}

	err = h.uc.ResetPassWord(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword.uc.ResetPassword: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

func (h handler) VerifyRequest(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processVerifyRequestRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.VerifyRequest.process: %v", err)
		response.Error(c, err)
		return
	}

	h.uc.VerifyRequest(ctx, sc.Email)

	response.OK(c, nil)

}
func (h handler) VerifyUser(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processVerifyUserRequesr(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.VerifyUser.uc.verifyuser: %v", err)
		response.Error(c, err)
		return
	}
	err = h.uc.VerifyUser(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.VerifyUser.uc.verifyuser: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)

}
func (h handler) DistributeNewToken(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processDistributeNewTokenRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.DistributeNewToken %v", err)
		response.Error(c, err)
		return
	}
	e, er := h.uc.DistributeNewToken(ctx, req.toInput())
	if er != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.DistributeNewToken.ToInput %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, h.newDistributeNewTokenResp(e))
}
