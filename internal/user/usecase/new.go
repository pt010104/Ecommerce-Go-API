package usecase

import (
	"github.com/pt010104/api-golang/internal/email"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUsecase struct {
	l         log.Logger
	repo      user.Repository
	emailUC   email.UseCase
	redisRepo user.Redis
	mediaUC   media.UseCase
}

func New(l log.Logger, repo user.Repository, emailUC email.UseCase, redisRepo user.Redis, mediaUC media.UseCase) user.UseCase {
	return &implUsecase{
		l:         l,
		repo:      repo,
		emailUC:   emailUC,
		redisRepo: redisRepo,
		mediaUC:   mediaUC,
	}
}
