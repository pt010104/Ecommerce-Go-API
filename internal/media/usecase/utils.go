package usecase

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUsecase) generateFilename(userID primitive.ObjectID) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s", timestamp, userID.Hex())
}

func (uc implUsecase) determineFolder(userID primitive.ObjectID, shopID primitive.ObjectID) string {
	if !shopID.IsZero() {
		return fmt.Sprintf("shops/%s", shopID.Hex())
	}
	return fmt.Sprintf("users/%s", userID.Hex())
}
