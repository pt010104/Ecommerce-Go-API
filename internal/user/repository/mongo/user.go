package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

func (repo implRepo) getUserCollection() mongo.Collection {
	return *repo.database.Collection(userCollection)
}

func (repo implRepo) CreateUser(ctx context.Context, opt user.CreateUserOption) (models.User, error) {
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

func (repo implRepo) GetUser(ctx context.Context, opt user.GetUserOption) (models.User, error) {
	col := repo.getUserCollection()

	filter, err := repo.buidUserQuery(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.GetUser.buidUserQuery: %v", err)
		return models.User{}, err
	}

	var user models.User

	err = col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		repo.l.Errorf(ctx, "survey.repository.mongo.GetUser.col.FindOne: %v", err)
		return models.User{}, err
	}

	return user, nil
}
func (repo implRepo) UpdateUser(ctx context.Context, opt user.UpdateUserOption) (models.User, error) {

	col := repo.getUserCollection()

	filter, err := repo.buildUserDetailQuery(ctx, opt.ID)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.GetUser.buildUserDetailQuery: %v", err)
		return models.User{}, err

	}

	update, nm, err := repo.buildUpdateUserModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.UpdateUser.buildUpdateUserModel: %v", err)
		return models.User{}, err
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.UpdateUser.col.UpdateOne: %v", err)
		return models.User{}, err
	}

	return nm, nil
}

func (repo implRepo) DetailUser(ctx context.Context, id string) (models.User, error) {
	col := repo.getUserCollection()

	filter, err := repo.buildUserDetailQuery(ctx, id)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.GetUser.buildUserDetailQuery: %v", err)
		return models.User{}, err
	}

	var user models.User

	err = col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		repo.l.Errorf(ctx, "survey.repository.mongo.Detail.col.FindOne: %v", err)
		return models.User{}, err
	}

	return user, nil
}
