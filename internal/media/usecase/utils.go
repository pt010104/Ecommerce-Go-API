package usecase

import (
	"fmt"
	"time"
)

func (uc implUsecase) generateFilename(userID string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s", timestamp, userID)
}

func (uc implUsecase) determineFolder(userID string, shopID string) string {
	if shopID != "" {
		return fmt.Sprintf("shops/%s", shopID)
	}
	return fmt.Sprintf("users/%s", userID)
}
