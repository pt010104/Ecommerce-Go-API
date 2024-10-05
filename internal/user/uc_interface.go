package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, input UseCaseType) (models.User, error)
	SignIn(ctx context.Context, input SignInType) (string, error)
}
