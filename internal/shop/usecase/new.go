package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/shop"
)

type implUsecase struct {
	l       log.Logger
	repo    shop.Repository
	adminUC admins.UseCase
}

var _ shop.UseCase = &implUsecase{}

func New(l log.Logger, repo shop.Repository, adminUC admins.UseCase) *implUsecase {
	return &implUsecase{
		l:       l,
		repo:    repo,
		adminUC: adminUC,
	}
}
func (uc *implUsecase) SetAdminUC(adminUC admins.UseCase) {
	uc.adminUC = adminUC
}
