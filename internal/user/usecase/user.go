package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
)

func (uc implUsecase) CreateUser(ctx context.Context, uct user.UseCaseType) (models.User, error) {
	hashedPass, err := uc.HashPassword(uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "error during  hashing pass : %v", err)
		return models.User{}, err
	}

	u, err := uc.repo.CreateUserRepo(ctx, user.RepoOption{

		Email:    uct.Email,
		Password: hashedPass,
		UserName: uct.UserName,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.CreateUser.repo.Create: %v", err)
		return models.User{}, err
	}
	return u, nil

}
