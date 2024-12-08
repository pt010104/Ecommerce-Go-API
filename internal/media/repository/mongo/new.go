package mongo

import (
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type implRepository struct {
	l        log.Logger
	database mongo.Database
}

func New(l log.Logger, database mongo.Database) media.Repository {
	return &implRepository{
		l:        l,
		database: database,
	}
}
