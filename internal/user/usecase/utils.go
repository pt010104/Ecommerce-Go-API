package usecase

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/mongo"

	"time"

	"golang.org/x/crypto/bcrypt"
)

func (uc implUsecase) hashPassword(password string) (string, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil

}

func (uc implUsecase) generateJWT(userName string, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userName,
		"iss": "US",

		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uc implUsecase) validateDataCreateUser(ctx context.Context, email string) error {
	_, err := uc.repo.GetUser(ctx, user.GetUserOption{
		Email: email,
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		uc.l.Errorf(ctx, "error during finding matching user: %v", err)
	}

	return user.ErrEmailExisted

}

func (uc implUsecase) generateVerificationJWT(userName string, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userName,
		"iss": "US",

		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"purpose": "email_verification",
	})
	tokenString, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
