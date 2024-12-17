package admins

import "errors"

var (
	ErrNoPermission = errors.New("you dont have permission to do this")
	ErrInvalidInput = errors.New("category must have name and description")
	ErrWrongBody    = errors.New("invalid object ID")
)
