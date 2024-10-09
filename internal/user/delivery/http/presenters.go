package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
)

type signupReq struct {
	UserName string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type forgetPasswordReq struct {
	Email string `json:"email" binding:"required"`
}
type signinReq struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	SessionID string `json:"session_id"`
}

func (r signupReq) toInput() user.CreateUserInput {
	return user.CreateUserInput{
		UserName: r.UserName,
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r signinReq) toInput() user.SignInType {
	return user.SignInType{
		Email:     r.Email,
		Password:  r.Password,
		SessionID: r.SessionID,
	}
}

type resetPasswordReq struct {
	UserID      string `json:"user_id" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	Token       string
}

func (r resetPasswordReq) toInput() user.ResetPasswordInput {
	return user.ResetPasswordInput{
		UserId:  r.UserID,
		NewPass: r.NewPassword,
		Token:   r.Token,
	}
}

type SignUpResponse struct {
	email    string
	username string
}

func ResponseSignUp(u models.User) SignUpResponse {
	return SignUpResponse{email: u.Email, username: u.UserName}
}

type detailResp struct {
	Email    string `json:"email"`
	UserName string `json:"name"`
}

func (h handler) newDetailResp(u models.User) detailResp {
	return detailResp{
		Email:    u.Email,
		UserName: u.UserName,
	}
}

type signInResp struct {
	SessionID string     `json:"session_id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Token     user.Token `json:"token"`
}

func (h handler) newSignInResp(output user.SignInOutput) signInResp {
	return signInResp{
		Email:     output.User.Email,
		Username:  output.User.UserName,
		Token:     output.Token,
		SessionID: output.SessionID,
	}
}
