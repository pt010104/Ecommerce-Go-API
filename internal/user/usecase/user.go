package usecase

import (
	"context"
	"time"

	"errors"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/mongo"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func (uc implUsecase) CreateUser(ctx context.Context, uct user.UseCaseType) (models.User, error) {
	err := uc.validateDataCreateUser(ctx, uct.Email)
	if err != nil {
		uc.l.Errorf(ctx, "error during validate data: %v", err)
		return models.User{}, err
	}

	hashedPass, err := uc.hashPassword(uct.Password)
	if err != nil {
		uc.l.Errorf(ctx, "error during  hashing pass : %v", err)
		return models.User{}, err
	}

	u, err := uc.repo.CreateUser(ctx, user.CreateUserOption{

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
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: sit.Email,
	})
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

	kt, err := uc.repo.CreateKeyToken(ctx, u.ID, sit.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "error during finding matching user: %v", err)
		return "", err
	}

	payload := jwt.Payload{
		UserID:    u.ID.Hex(),
		Refresh:   false,
		SessionID: sit.SessionID,
	}

	expirationTime := time.Hour * 24
	token, err := jwt.Sign(payload, expirationTime, kt.SecretKey)
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return "", err
	}
	uc.repo.DeleteRecord(ctx, u.ID.Hex(), sit.SessionID)
	return token, nil
}
func (uc implUsecase) ForgetPasswordRequest(ctx context.Context, email string) (token string, err error) {
	u, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: email,
	})
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.ForgetPasswordRequest: %v", err)
		return "", err
	}
	payload := jwt.Payload{
		UserID:  u.ID.Hex(),
		Refresh: false,
		Type:    "reset-request",
	}
	expirationTime := time.Hour * 1
	token, err = jwt.Sign(payload, expirationTime, os.Getenv("SUPER_SECRET_KEY"))
	if err != nil {
		uc.l.Errorf(ctx, "error signing token: %v", err)
		return "", err
	}
	err1 := uc.emailUC.SendVerificationEmail(u.Email, token)
	if err1 != nil {
		uc.l.Errorf(ctx, "user.usecase.ForgetPasswordRequest: %v", err)
	}
	return token, nil

}
func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (models.User, error) {
	u, err := uc.repo.DetailUser(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			uc.l.Warnf(ctx, "user.usecase.Detail.repo.DetailUser: %v", err)
			return models.User{}, user.ErrUserNotFound
		}
		uc.l.Errorf(ctx, "user.usecase.Detail.repo.DetailUser: %v", err)
		return models.User{}, err
	}
	uc.repo.DeleteRecord(ctx, sc.UserID, sc.SessionID)
	return u, nil
}
func (uc implUsecase) LogOut(ctx context.Context, sc models.Scope) {
	err := uc.repo.DeleteRecord(ctx, sc.UserID, sc.SessionID)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.LogOut.repo.DeleteRecord")
	}

}
