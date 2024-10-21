package repository

import (
	"log"

	"github.com/pt010104/api-golang/internal/product"
)

type implRepo struct {
	l  log.Logger
	uc product.UseCase
}

func New(l log.Logger, uc product.UseCase) implRepo {
	return implRepo{
		l:  l,
		uc: uc,
	}
}
