package mongo

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type implRepo struct {
	l        log.Logger
	database mongo.Database
}

func New(l log.Logger, database mongo.Database) cart.Repo {
	return &implRepo{
		l:        l,
		database: database,
	}
}
