package usecase

import (
	"fmt"

	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) generateFilename(userID string) string {
	timestamp := util.Now().UnixNano()
	return fmt.Sprintf("%d_%s", timestamp, userID)
}

func (uc implUsecase) determineFolder(userID string, shopID string) string {
	if shopID != "" {
		return fmt.Sprintf("shops/%s", shopID)
	}
	return fmt.Sprintf("users/%s", userID)
}
