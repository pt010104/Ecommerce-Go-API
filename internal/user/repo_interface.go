package user

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockery --name=Repository
type Repository interface {
	CreateUser(context context.Context, opt CreateUserOption) (models.User, error)
	GetUser(ctx context.Context, opt GetUserOption) (models.User, error)
	DetailUser(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, opt UpdateUserOption) (models.User, error)
	ListUser(ctx context.Context, opt ListUserOption) ([]models.User, error)

	CreateRequestToken(ctx context.Context, id primitive.ObjectID, token string) (models.RequestToken, error)
	UpdatePatchUser(ctx context.Context, opt UpdateUserOption) (models.User, error)
	UpdateRequestToken(ctx context.Context, opt UpdateRequestTokenOption) error
	DetailRequestToken(ctx context.Context, JWT string) (models.RequestToken, error)
	UpdateKeyToken(ctx context.Context, opt UpdateKeyTokenInput) error
	DetailKeyToken(ctx context.Context, userID string, sessionID string) (models.KeyToken, error)
	CreateKeyToken(context context.Context, opt CreateKeyTokenOption) (models.KeyToken, error)
	DeleteKeyToken(context context.Context, userID string, sessionID string) error
}
