package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repo interface {
	CreateUserRepo(context context.Context, opt RepoOption) (models.User, error)
	GetUserRepo(ctx context.Context, opt string) (models.User, error)
}
