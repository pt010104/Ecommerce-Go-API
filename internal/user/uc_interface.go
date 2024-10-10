package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	SignIn(ctx context.Context, input SignInType) (SignInOutput, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.User, error)
	LogOut(ctx context.Context, sc models.Scope)

	ForgetPasswordRequest(ctx context.Context, email string) (token string, err error)
	VerifyRequest(ctx context.Context, email string) (token string, err error)
	ResetPassWord(ctx context.Context, input ResetPasswordInput) error
	VerifyUser(ctx context.Context, input VerifyUserInput) error
}
