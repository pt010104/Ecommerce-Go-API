package user

import "github.com/pt010104/api-golang/internal/models"

type CreateUserInput struct {
	Name     string
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
type VerifyEmailInput struct {
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
type DistributeNewTokenInput struct {
	UserId       string
	SessionID    string
	RefreshToken string
}
type DistributeNewTokenOutput struct {
	Token Token
}
type UpdateInput struct {
	Email   string
	Name    string
	MediaID string
}
type DetailUserOutput struct {
	User       models.User
	Avatar_URL string
}

type AddAddressInput struct {
	Street   string
	District string
	City     string
	Province string
	Phone    string
	Default  bool
}

type UpdateAddressInput struct {
	AddressID string
	Street    string
	District  string
	City      string
	Province  string
	Phone     string
	Default   bool
}

type DetailAddressOutput struct {
	Addressess []models.Address
}
