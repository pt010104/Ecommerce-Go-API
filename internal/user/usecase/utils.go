package usecase

import (
	"context"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pt010104/api-golang/internal/user"

	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/pt010104/api-golang/pkg/util"
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
func (uc implUsecase) validateDataUser(ctx context.Context, email, password string) error {
	// Validate email
	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return user.ErrInvalidEmailFormat
	}

	// Validate password
	passwordRegex := `^[A-Za-z\d]{8,}$`
	hasLetter := false
	hasDigit := false

	for _, ch := range password {
		if ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') {
			hasLetter = true
		}
		if '0' <= ch && ch <= '9' {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit || !regexp.MustCompile(passwordRegex).MatchString(password) {
		return user.ErrInvalidPasswordFormat
	}

	return nil
}

func (uc implUsecase) createKeyToken(ctx context.Context) (string, string, string, error) {
	sessionID, err := util.GenerateRandomString(32)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Signin.createKeyToken: %v", err)
		return "", "", "", err
	}

	secretKey, err := util.GenerateRandomString(32)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Signin.createKeyToken: %v", err)
		return "", "", "", err
	}

	refreshToken, err := util.GenerateRandomString(32)
	if err != nil {
		uc.l.Errorf(ctx, "user.usecase.Signin.createKeyToken: %v", err)
		return "", "", "", err
	}

	return sessionID, secretKey, refreshToken, nil
}
