package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/vouchers"
)

type implUsecase struct {
	repo   vouchers.Repository
	l      log.Logger
	shopUc shop.UseCase
}

var _ vouchers.UseCase = implUsecase{}

func New(repo vouchers.Repository, l log.Logger, shopUC shop.UseCase) implUsecase {
	return implUsecase{
		repo:   repo,
		l:      l,
		shopUc: shopUC,
	}
}
