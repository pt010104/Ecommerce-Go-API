package media

import "errors"

var (
	ErrRequireField  = errors.New("require field")
	ErrMediaFailed   = errors.New("media failed")
	ErrMediaPending  = errors.New("media pending")
	ErrInvalidStatus = errors.New("invalid status")
)
