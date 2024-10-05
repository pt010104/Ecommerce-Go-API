package usecase

import (
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
	"time"
)

func (uc implUsecase) HashPassword(password string) (string, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil

}

func (uc implUsecase) GenerateJWT(userName string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userName,
		"iss": "US",

		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := claims.SignedString([]byte("supersecret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
