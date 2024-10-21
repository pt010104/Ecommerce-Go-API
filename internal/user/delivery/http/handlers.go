package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary User Sign Up
// @Schemes http https
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Name body string true "Name"
// @Param Email body string true "Email"
// @Param Password body string true "Password"
//
// @Success 200 {object} signUpResponse
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/sign-up [POST]
func (h handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := h.processSignupRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.SignUp.processSignupRequest: %v", err)
		response.Error(c, err)
		return
	}

	u, err := h.uc.CreateUser(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.SignUp.uc.CreateUser: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newSignUpResponse(u))
}

// @Summary User Sign In
// @Schemes http https
// @Description Authenticate user and provide a JWT token
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param session-id header string false "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param email body string true "Email"
// @Param password body string true "Password"
//
// @Success 200 {object} signInResp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/sign-in [POST]
func (h handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processSignInRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.SignIn.processSignInRequest: %v", err)
		response.Error(c, err)
		return
	}

	output, err := h.uc.SignIn(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.SignIn.uc.SignIn: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	signInResp := h.newSignInResp(output)
	response.OK(c, signInResp)
}

// @Summary Forget Password Request
// @Schemes http https
// @Description Initiate a password reset request
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param email body string true "Email"
//
// @Success 200 {object} interface{}
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/forget-password [POST]
func (h handler) ForgetPasswordRequest(c *gin.Context) {
	ctx := c.Request.Context()

	sc, err := h.processForgetPasswordRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.ForgetPasswordRequest.processForgetPasswordRequest: %v", err)
		response.Error(c, err)
		return
	}

	h.uc.ForgetPasswordRequest(ctx, sc.Email)

	response.OK(c, nil)
}

// @Summary User Sign Out
// @Schemes http https
// @Description Sign out the current user
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
//
// @Success 200 {object} interface{}
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/signout [POST]
func (h handler) SignOut(c *gin.Context) {
	ctx := c.Request.Context()
	sc, err := h.processLogOutRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.SignOut.processLogOutRequest: %v", err)
		response.Error(c, err)
		return
	}
	err = h.uc.LogOut(ctx, sc)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.SignOut.uc.LogOut: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// @Summary Get User Details
// @Schemes http https
// @Description Retrieve detailed information about a specific user
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param id path string true "User ID"
//
// @Success 200 {object} detailResp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 404 {object} response.Resp "User Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/{id} [GET]
func (h handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	id, sc, err := h.processDetailRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Detail.processDetailRequest: %v", err)
		response.Error(c, err)
		return
	}

	u, err := h.uc.Detail(ctx, sc, id)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.Detail.uc.Detail: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newDetailResp(u))
}

// @Summary Reset Password
// @Schemes http https
// @Description Reset user's password using a valid reset token
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param token query string true "Reset token"
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param password body string true "New password"
//
// @Success 200 {object} interface{}
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/reset-password [POST]
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

// @Summary Verify Email
// @Schemes http https
// @Description Verify the validity of a password reset request
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param email body string true "Email"
//
// @Success 200 {object} interface{}
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 404 {object} response.Resp "Request Not Found"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/verify-request [POST]
func (h handler) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := h.processVerifyEmailRequest(c)
	if err != nil {
		h.l.Error(ctx, "user.delivery.http.handler.VerifyEmail.processVerifyEmailRequest: %v", err)
		response.Error(c, err)
		return
	}

	h.uc.VerifyEmail(ctx, req.Email)

	response.OK(c, nil)
}

// @Summary Verify User Account
// @Schemes http https
// @Description Complete user account verification
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param token query string true "Verification token"
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
//
// @Success 200 {object} interface{}
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/verify-user [POST]
func (h handler) VerifyUser(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processVerifyUserRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.VerifyUser.uc.process: %v", err)
		response.Error(c, err)
		return
	}
	err = h.uc.VerifyUser(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.VerifyUser.uc.VerifyUser: %v", err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// @Summary Distribute New Token
// @Schemes http https
// @Description Generate and distribute a new JWT token
// @Tags User
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param token query string true "Verification token"
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
//
// @Success 200 {object} distributeNewTokenResp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 401 {object} response.Resp "Unauthorized"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/users/distribute-token [POST]
func (h handler) DistributeNewToken(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := h.processDistributeNewTokenRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.DistributeNewToken.processDistributeNewTokenRequest: %v", err)
		response.Error(c, err)
		return
	}
	e, er := h.uc.DistributeNewToken(ctx, req.toInput())
	if er != nil {
		h.l.Errorf(ctx, "user.delivery.http.handler.DistributeNewToken.DistributeNewToken: %v", er)
		response.Error(c, er)
		return
	}
	response.OK(c, h.newDistributeNewTokenResp(e))
}
