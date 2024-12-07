package mongo

import (
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/log"

	"go.mongodb.org/mongo-driver/mongo"
)

type implRepo struct {
	l        log.Logger
	database mongo.Database
}

var _ vouchers.Repository = implRepo{}

func New(l log.Logger, database mongo.Database) implRepo {
	return implRepo{
		l:        l,
		database: database,
	}
}
