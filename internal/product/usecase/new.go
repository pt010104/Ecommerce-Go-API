package usecase

import (
	"log"

	"github.com/pt010104/api-golang/internal/product"
)

type implUsecase struct {
	l    log.Logger
	repo product.Repo
}

func New(l log.Logger, repo product.Repo) implUsecase {
	return implUsecase{
		l:    l,
		repo: repo,
	}
}
