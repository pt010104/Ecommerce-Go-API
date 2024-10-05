package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	CreateUserRepo(context context.Context, opt RepoOption) (models.User, error)
	GetUserRepo(ctx context.Context, opt string) (models.User, error)
	CreateKeyToken(context context.Context, UserId primitive.ObjectID) (models.KeyToken, error)
}
