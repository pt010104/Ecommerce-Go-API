package http

import (
	"errors"
	"regexp"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
)

var (
	errInvalidEmail    = errors.New("Invalid email format")
	errInvalidPassword = errors.New("Password must be at least 8 characters long and include letters and numbers")
	errInvalidUserName = errors.New("Username must be at least 3 characters long")

	errUserNameExisted = errors.New("Your name must be unique")
)

var passwordRegex = `^[A-Za-z\d]{8,}$`

func validateEmail(email string) bool {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func validatePassword(password string) bool {
	hasLetter := false
	hasDigit := false

	for _, ch := range password {
		if ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') {
			hasLetter = true
		}
		if '0' <= ch && ch <= '9' {
			hasDigit = true
		}
	}

	if hasLetter && hasDigit {
		return regexp.MustCompile(passwordRegex).MatchString(password)
	}

	return false
}

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

func (r forgetPasswordReq) toInput() user.ForgetPasswordRequest {
	return user.ForgetPasswordRequest{
		Email: r.Email,
	}
}

func (r signupReq) validate() error {

	if len(r.UserName) < 3 {
		return errInvalidUserName
	}

	if !validateEmail(r.Email) {
		return errInvalidEmail
	}

	if !validatePassword(r.Password) {
		return errInvalidPassword
	}

	return nil
}

func (r signinReq) validate() error {

	if !validateEmail(r.Email) {
		return errInvalidEmail
	}

	if !validatePassword(r.Password) {
		return errInvalidPassword
	}

	if r.SessionID == "" {
		return errWrongBody
	}

	return nil
}
func (r resetPasswordReq) toInput() user.ResetPassWordReq {
	return user.ResetPassWordReq{
		NewPassword: r.NewPassword,
	}
}
func (r signupReq) toInput() user.UseCaseType {
	return user.UseCaseType{
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
	NewPassword string
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
	Session  string `json: "session_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Jwt      string `json:"token"`
}

func (h handler) newSignInResp(u models.User, session string, jwt string) signInResp {
	return signInResp{
		Email:    u.Email,
		Username: u.UserName,
		Jwt:      jwt,
		Session:  session,
	}
}
