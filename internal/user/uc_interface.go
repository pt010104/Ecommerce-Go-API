package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, input UseCaseType) (models.User, error)
	SignIn(ctx context.Context, input SignInType) (string, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.User, error)
	LogOut(ctx context.Context, sc models.Scope)
	ForgetPasswordRequest(ctx context.Context, email string) (token string, err error)
}
