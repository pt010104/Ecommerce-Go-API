package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

func (uc implUsecase) HashPassword(password string) (string, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil

}
