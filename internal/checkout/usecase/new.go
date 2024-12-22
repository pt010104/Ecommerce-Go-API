package usecase

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/checkout"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUseCase struct {
	l      log.Logger
	repo   checkout.Repo
	shopUC shop.UseCase
	cartUC cart.UseCase
}

func New(l log.Logger, repo checkout.Repo, shopUC shop.UseCase, cartUC cart.UseCase) checkout.UseCase {
	return &implUseCase{
		l:      l,
		repo:   repo,
		shopUC: shopUC,
		cartUC: cartUC,
	}
}
