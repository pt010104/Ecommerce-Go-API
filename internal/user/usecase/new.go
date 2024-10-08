package usecase

import (
	"github.com/pt010104/api-golang/internal/email"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUsecase struct {
	l       log.Logger
	repo    user.Repo
	emailUC email.UseCase
}

func New(l log.Logger, repo user.Repo, emailUC email.UseCase) implUsecase {
	return implUsecase{
		l:       l,
		repo:    repo,
		emailUC: emailUC,
	}
}
