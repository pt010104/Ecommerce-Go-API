package usecase

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUseCase struct {
	l         log.Logger
	redisRepo order.Redis
	repo      order.Repo
	shopUC    shop.UseCase
	cartUC    cart.UseCase
}

func New(l log.Logger, repo order.Repo, shopUC shop.UseCase, cartUC cart.UseCase, redisRepo order.Redis) order.UseCase {
	return &implUseCase{
		l:         l,
		repo:      repo,
		shopUC:    shopUC,
		cartUC:    cartUC,
		redisRepo: redisRepo,
	}
}
