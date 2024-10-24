package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	"github.com/pt010104/api-golang/internal/admin"
)

type implUsecase struct {
	repo admin.Repo
	l    log.Logger
}

func New(repo admin.Repo, l log.Logger) implUsecase {
	return implUsecase{
		repo: repo,
		l:    l,
	}
}
