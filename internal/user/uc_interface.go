package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	SignIn(ctx context.Context, input SignInType) (SignInOutput, error)
	Detail(ctx context.Context, sc models.Scope, id string) (DetailUserOutput, error)
	LogOut(ctx context.Context, sc models.Scope) error
	GetModel(ctx context.Context, id string) (models.User, error)
	VerifyUser(ctx context.Context, input VerifyUserInput) error
	Update(ctx context.Context, sc models.Scope, input UpdateInput) (DetailUserOutput, error)

	DetailKeyToken(ctx context.Context, userID string, sessionID string) (models.KeyToken, error)
	ForgetPasswordRequest(ctx context.Context, email string) (token string, err error)
	VerifyEmail(ctx context.Context, email string) (token string, err error)
	ResetPassWord(ctx context.Context, input ResetPasswordInput) error
	DistributeNewToken(ctx context.Context, input DistributeNewTokenInput) (output DistributeNewTokenOutput, er error)
}
