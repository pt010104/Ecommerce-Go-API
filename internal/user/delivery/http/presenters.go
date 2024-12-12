package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/mongo"
)

type signupReq struct {
	Name     string `json:"name" binding:"required"`
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
type verifyRequestReq struct {
	Email string `json:"email" binding:"required"`
}

func (r signupReq) toInput() user.CreateUserInput {
	return user.CreateUserInput{
		Name:     r.Name,
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

type verifyUserReq struct {
	UserID string
	Token  string
}

func (r verifyUserReq) validate() error {
	if r.UserID == "" {
		return errWrongHeader
	}

	if r.Token == "" {
		return errWrongQuery
	}

	return nil
}

type resetPasswordReq struct {
	UserID      string
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
func (r verifyUserReq) toInput() user.VerifyUserInput {
	return user.VerifyUserInput{
		UserId: r.UserID,
		Token:  r.Token}
}

type signUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h handler) newSignUpResponse(u models.User) signUpResponse {
	return signUpResponse{
		ID:    u.ID.Hex(),
		Email: u.Email,
		Name:  u.Name,
	}
}

type detailResp struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`

	Avatar Avatar_obj `json:"avatar"`
}

func (h handler) newDetailResp(u user.DetailUserOutput) detailResp {
	return detailResp{
		ID:    u.User.ID.Hex(),
		Email: u.User.Email,
		Name:  u.User.Name,
		Avatar: Avatar_obj{
			MediaID: u.User.MediaID.Hex(),
			URL:     u.Avatar_URL,
		},
	}
}

type distributeNewTokenReq struct {
	UserId       string
	SessionID    string
	RefreshToken string
}

func (r distributeNewTokenReq) validate() error {
	if r.UserId == "" || r.SessionID == "" || r.RefreshToken == "" {
		return errWrongHeader
	}

	return nil
}

func (r distributeNewTokenReq) toInput() user.DistributeNewTokenInput {
	return user.DistributeNewTokenInput{
		UserId:       r.UserId,
		SessionID:    r.SessionID,
		RefreshToken: r.RefreshToken,
	}
}

type Avatar_obj struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
type signInResp struct {
	ID        string     `json:"id"`
	SessionID string     `json:"session_id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Token     user.Token `json:"token"`
}

func (h handler) newSignInResp(output user.SignInOutput) signInResp {
	return signInResp{
		ID:        output.User.ID.Hex(),
		Email:     output.User.Email,
		Username:  output.User.Name,
		Token:     output.Token,
		SessionID: output.SessionID,
	}
}

type distributeNewTokenResp struct {
	NewAccessToken  string `json:"new_access_token"`
	NewRefreshToken string `json:"new_refresh_token"`
}

func (h handler) newDistributeNewTokenResp(output user.DistributeNewTokenOutput) distributeNewTokenResp {
	return distributeNewTokenResp{
		NewAccessToken:  output.Token.AccessToken,
		NewRefreshToken: output.Token.RefreshToken,
	}
}

type UpdateAvatarReq struct {
	MediaID string `json:"media_id" binding:"required"`
}

func (r UpdateAvatarReq) toInput() user.UpdateAvatarInput {
	return user.UpdateAvatarInput{
		MediaID: r.MediaID,
	}
}
func (r UpdateAvatarReq) validate() error {
	if r.MediaID == "" || !mongo.IsObjectID(r.MediaID) {
		return errWrongBody
	}

	return nil
}

type UpdateAvatarResp struct {
	MediaID string `json:"media_id"`
}

func (h handler) newUpdateAvatarResp(user models.User) UpdateAvatarResp {
	return UpdateAvatarResp{
		MediaID: user.MediaID.Hex(),
	}
}
