package usecase

import "github.com/pt010104/api-golang/internal/models"

var validStatus = []string{
	models.MediaStatusPending,
	models.MediaStatusUploaded,
	models.MediaStatusFailed,
	models.MediaStatusDrafted,
}
