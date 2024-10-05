package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection     = "users"
	keyTokenCollection = "tokens"
)

func (repo implRepo) getUserCollection() mongo.Collection {
	return *repo.database.Collection(userCollection)
}
func (repo implRepo) getKeyTokenCollection() mongo.Collection {
	return *repo.database.Collection(keyTokenCollection)
}

func (repo implRepo) CreateUserRepo(ctx context.Context, opt user.RepoOption) (models.User, error) {
	col := repo.getUserCollection()

	u, err := repo.buildUserModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.buldUserModel: %v", err)
		return models.User{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.InsertOne: %v", err)
		return models.User{}, err
	}
	return u, nil
}
func (repo implRepo) GetUserRepo(ctx context.Context, email string) (models.User, error) {
	col := repo.getUserCollection()

	filter := bson.M{"email": email}

	var user models.User

	err := col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}
func (repo implRepo) CreateKeyToken(ctx context.Context, userId primitive.ObjectID) (models.KeyToken, error) {
	col := repo.getKeyTokenCollection()
	u, err := repo.buildKeyTokenModel(ctx, userId)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.buldUserModel: %v", err)
		return models.KeyToken{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.InsertOne: %v", err)
		return models.KeyToken{}, err
	}
	return u, nil

}
