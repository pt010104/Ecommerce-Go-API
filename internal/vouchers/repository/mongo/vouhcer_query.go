package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (repo implRepo) buildVoucherDetailQuery(ctx context.Context, sc models.Scope, code string) bson.M {
	filter := bson.M{}

	filter["code"] = code

	return filter

}
