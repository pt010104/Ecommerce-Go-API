package mongo

import (
	"github.com/pt010104/api-golang/internal/checkout"
	"github.com/pt010104/api-golang/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type implRepo struct {
	l        log.Logger
	database mongo.Database
}

func New(l log.Logger, database mongo.Database) checkout.Repo {
	return &implRepo{
		l:        l,
		database: database,
	}
}
