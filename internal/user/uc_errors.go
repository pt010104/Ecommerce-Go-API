package user

import "errors"

var (
	ErrEmailExisted = errors.New("email has already been registered")
	ErrUserNotFound = errors.New("user not found")
)
