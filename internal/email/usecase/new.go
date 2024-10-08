package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"
)

type implUsecase struct {
	l log.Logger
}

func New(l log.Logger) implUsecase {
	return implUsecase{
		l: l,
	}
}
