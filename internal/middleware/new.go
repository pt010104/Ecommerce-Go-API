package middleware

import (
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type Middleware struct {
	l      log.Logger
	shopUC shop.UseCase
	userUC user.UseCase
}

func New(l log.Logger, shopUC shop.UseCase, userUC user.UseCase) Middleware {
	return Middleware{
		l:      l,
		shopUC: shopUC,
		userUC: userUC,
	}
}
