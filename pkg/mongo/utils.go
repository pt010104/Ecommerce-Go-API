package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDFromHexOrNil returns an ObjectID from the provided hex representation.
func ObjectIDFromHexOrNil(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}

func ObjectIDsFromHexOrNil(ids []string) []primitive.ObjectID {
	objIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objIDs[i] = ObjectIDFromHexOrNil(id)
	}
	return objIDs
}

func HexFromObjectIDOrNil(id primitive.ObjectID) string {
	return id.Hex()
}

func HexFromObjectIDsOrNil(ids []primitive.ObjectID) []string {
	hexIds := make([]string, len(ids))
	for i, id := range ids {
		hexIds[i] = HexFromObjectIDOrNil(id)
	}
	return hexIds
}

func BuildQueryWithSoftDelete(query bson.M) bson.M {
	query["deleted_at"] = nil
	return query
}
func BuildShopScopeQuery(ctx context.Context, l log.Logger, sc models.Scope) (bson.M, error) {
	filter := bson.M{}

	if sc.ShopID != "" {
		ShopID, err := primitive.ObjectIDFromHex(sc.ShopID)
		if err != nil {
			return nil, err
		}
		filter["shop_id"] = ShopID
	}

	return filter, nil
}

func BuildScopeQuery(ctx context.Context, l log.Logger, sc models.Scope) (bson.M, error) {
	filter := bson.M{}

	if sc.UserID != "" {
		UserID, err := primitive.ObjectIDFromHex(sc.UserID)
		if err != nil {
			return nil, err
		}
		filter["user_id"] = UserID
	}

	return filter, nil
}

func IsObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

func ObjectIDsFromHex(ids []string) ([]primitive.ObjectID, error) {
	objIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs[i] = objID
	}

	return objIDs, nil
}
