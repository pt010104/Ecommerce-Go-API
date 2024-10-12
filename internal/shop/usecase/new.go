package usecase

import (
	"log"

	"github.com/pt010104/api-golang/internal/shop"
)

type implUsecase struct {
	l    log.Logger
	repo shop.Repo
}

func New(l log.Logger, repo shop.Repo) implUsecase {
	return implUsecase{
		l:    l,
		repo: repo,
	}
}
