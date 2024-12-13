package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/user"
)

type implUsecase struct {
	l       log.Logger
	repo    shop.Repository
	adminUC admins.UseCase
	userUC  user.UseCase
}

var _ shop.UseCase = &implUsecase{}

func New(l log.Logger, repo shop.Repository, adminUC admins.UseCase, userUC user.UseCase) *implUsecase {
	return &implUsecase{
		l:       l,
		repo:    repo,
		adminUC: adminUC,
		userUC:  userUC,
	}
}
func (uc *implUsecase) SetAdminUC(adminUC admins.UseCase) {
	uc.adminUC = adminUC
}
