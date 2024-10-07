package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	CreateUser(context context.Context, opt CreateUserOption) (models.User, error)
	GetUser(ctx context.Context, opt GetUserOption) (models.User, error)
	DetailUser(ctx context.Context, id string) (models.User, error)

	DetailKeyToken(ctx context.Context, userID string, sessionID string) (models.KeyToken, error)
	CreateKeyToken(context context.Context, UserId primitive.ObjectID) (models.KeyToken, error)
}
