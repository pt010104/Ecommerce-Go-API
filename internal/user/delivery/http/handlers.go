package http

import (
	"errors"
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

	token, err := h.uc.SignIn(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Signin.uc.signIn: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, token)

}
func (h handler) ForgetPasswordRequest(c *gin.Context) {
	ctx := c.Request.Context()
	sc, err := h.processForgetPasswordRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.Signup.processForgetPassRequest: %v", err)
		response.Error(c, err)
		return
	}
	response.OK(c, "password reset request sent")
	h.uc.ForgetPasswordRequest(ctx, sc.Email)
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
	token := c.Query("token")
	if token == "" {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword: missing token in request")
		response.Error(c, errors.New("token is required"))
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword: missing userID in context")
		response.Error(c, errors.New("user ID not found"))
		return
	}

	var req struct {
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword.ShouldBindJSON: %v", err)
		response.Error(c, errors.New("new password is required"))
		return
	}
	if validatePassword(req.NewPassword) == false {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword.ShouldBindJSON: %v", "password should contain numbers and digits")
		response.Error(c, errors.New("wrong body"))
		return
	}
	valid, err1 := h.uc.IsJWTresetVaLID(ctx, token)
	if err1 != nil {

		return
	}

	if valid == false {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword0: %v", "token is used")
		response.Error(c, errors.New("token is used"))
		return
	}
	err := h.uc.ResetPassWord(ctx, userID.(string), req.NewPassword)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword.uc.ResetPassword: %v", err)
		response.Error(c, errors.New("failed to reset password"))
		return
	}
	err = h.uc.MartJWTasUsed(ctx, token)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.ResetPassword.uc.Resetpass.MarkJWTasUsed: %v", err)
	}

	response.OK(c, gin.H{"message": "Password reset successfully"})
}
