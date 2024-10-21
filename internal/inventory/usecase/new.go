package usecase

import (
	"github.com/pt010104/api-golang/internal/inventory"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUseCase struct {
	l    log.Logger
	repo inventory.Repo
}

func New(l log.Logger, repo inventory.Repo) implUseCase {
	return implUseCase{
		l:    l,
		repo: repo,
	}
}
