package usecase

import (
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUsecase struct {
	l    log.Logger
	repo user.Repo
}

func New(l log.Logger, repo user.Repo) implUsecase {
	return implUsecase{
		l:    l,
		repo: repo,
	}
}
