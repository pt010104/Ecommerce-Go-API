package media

import (
	"github.com/pt010104/api-golang/internal/models"
)

type UploadOption struct {
	Folder   string
	FileName string
}

type UpdateOption struct {
	Model  models.Media
	Status string
	URL    string
}
