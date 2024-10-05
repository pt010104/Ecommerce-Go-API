package usecase

import (
	"context"

	"errors"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"golang.org/x/crypto/bcrypt"
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

func (uc implUsecase) SignIn(ctx context.Context, sit user.SignInType) (string, error) {
	u, err := uc.repo.GetUserRepo(ctx, sit.Email)
	if err != nil {
		uc.l.Errorf(ctx, "error during finding matching user: %v", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(sit.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			uc.l.Warnf(ctx, "password mismatch for user: %v", u.Email)
			return "", errors.New("invalid email or password")
		}
		uc.l.Errorf(ctx, "error comparing passwords: %v", err)
		return "", err
	}
	token, err := uc.GenerateJWT(u.UserName)
	if err != nil {
		uc.l.Errorf(ctx, "error generating JWT: %v", err)
		return "", err
	}

	return token, nil

}
