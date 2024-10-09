package middleware

import (
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type Middleware struct {
	l    log.Logger
	repo user.Repo
}

func New(l log.Logger, repo user.Repo) Middleware {
	return Middleware{
		l:    l,
		repo: repo,
	}
}
