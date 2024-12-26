package usecase

import (
	"github.com/pt010104/api-golang/internal/email"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUsecase struct {
	l log.Logger
}

func New(l log.Logger) email.UseCase {
	return implUsecase{
		l: l,
	}
}
