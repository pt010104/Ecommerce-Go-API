package usecase

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUseCase struct {
	l      log.Logger
	repo   cart.Repo
	shopUc shop.UseCase
}

func New(l log.Logger, repo cart.Repo, shopUc shop.UseCase) implUseCase {
	return implUseCase{
		l:      l,
		repo:   repo,
		shopUc: shopUc,
	}
}
