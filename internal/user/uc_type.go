package user

import "github.com/pt010104/api-golang/internal/models"

type CreateUserInput struct {
	UserName string
	Password string
	Email    string
}

type SignInType struct {
	Email     string
	Password  string
	SessionID string
}

type ForgetPasswordRequest struct {
	Email string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
type SignInOutput struct {
	User      models.User
	Token     Token
	SessionID string
}
type VerifyRequestInput struct {
	Email string
}
type ResetPasswordInput struct {
	UserId  string
	NewPass string
	Token   string
}
type VerifyUserInput struct {
	UserId string
	Token  string
}
