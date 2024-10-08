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
	CreateRequestToken(ctx context.Context, id primitive.ObjectID, token string) (models.RequestToken, error)
	DetailKeyToken(ctx context.Context, userID string, sessionID string) (models.KeyToken, error)
	CreateKeyToken(context context.Context, UserId primitive.ObjectID, sessionID string) (models.KeyToken, error)
	DeleteRecord(context context.Context, userID string, sessionID string) error
}
