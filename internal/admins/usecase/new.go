package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/shop"
)

type implUsecase struct {
	repo   admins.Repository
	l      log.Logger
	shopUc shop.UseCase
}

func New(repo admins.Repository, l log.Logger, shopUC shop.UseCase) implUsecase {
	return implUsecase{
		repo:   repo,
		l:      l,
		shopUc: shopUC,
	}
}
