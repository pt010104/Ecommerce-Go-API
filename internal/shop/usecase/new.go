package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	"github.com/pt010104/api-golang/internal/shop"
)

type implUsecase struct {
	l    log.Logger
	repo shop.Repository
}

var _ shop.UseCase = implUsecase{}

func New(l log.Logger, repo shop.Repository) implUsecase {
	return implUsecase{
		l:    l,
		repo: repo,
	}
}
