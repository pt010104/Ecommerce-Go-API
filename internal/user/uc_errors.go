package user

import "errors"

var (
	ErrEmailExisted = errors.New("email has already been registered")
	ErrUserNotFound = errors.New("user not found")
	ErrTokenUsed    = errors.New("token is already used")

	ErrInvalidPasswordFormat = errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number and one special character")
	ErrInvalidEmailFormat    = errors.New("invalid email format")

	ErrInvalidUserName = errors.New("Username must be at least 3 characters long")

	ErrUserNotVerified = errors.New("user not verified")
)
