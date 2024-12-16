package shop

import "errors"

var (
	ErrInvalidInput                = errors.New("invalid input")
	ErrInvalidPhone                = errors.New("invalid phone")
	ErrShopExist                   = errors.New("shop exist")
	ErrShopDoesNotExist            = errors.New("shop doesnot exist")
	ErrNoPermissionToUpdate        = errors.New("no permission to update")
	ErrNonExistCategory            = errors.New("category doesnt exist")
	ErrNoPermissionToDeleteProduct = errors.New("no permission to delete product")
	ErrProductNotFound             = errors.New("can find product in this shop")
	ErrRequireField                = errors.New("require field")
)
