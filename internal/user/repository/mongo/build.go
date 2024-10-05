package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildUserModel(context context.Context, opt user.RepoOption) (models.User, error) {
	u := models.User{
		ID:       primitive.NewObjectID(),
		Email:    opt.Email,
		UserName: opt.UserName,
		Password: opt.Password,
	}
	return u, nil

}
