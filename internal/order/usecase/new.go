package usecase

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/email"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/internal/order/delivery/rabbitmq/producer"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUseCase struct {
	l         log.Logger
	redisRepo order.Redis
	repo      order.Repo
	shopUC    shop.UseCase
	cartUC    cart.UseCase
	prod      producer.Producer
	emailUC   email.UseCase
	userUC    user.UseCase
}

func New(l log.Logger, repo order.Repo, shopUC shop.UseCase, cartUC cart.UseCase, redisRepo order.Redis, prod producer.Producer, emailUC email.UseCase, userUC user.UseCase) order.UseCase {
	return &implUseCase{
		l:         l,
		repo:      repo,
		shopUC:    shopUC,
		cartUC:    cartUC,
		redisRepo: redisRepo,
		prod:      prod,
		emailUC:   emailUC,
		userUC:    userUC,
	}
}
