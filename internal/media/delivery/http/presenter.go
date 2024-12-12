package http

import (
	"mime/multipart"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type uploadRequest struct {
	Files []*multipart.FileHeader
}

func (r uploadRequest) toInput() media.UploadInput {
	files := make([][]byte, 0, len(r.Files))

	for _, file := range r.Files {
		f, err := file.Open()
		if err != nil {
			continue
		}
		defer f.Close()

		buf := make([]byte, file.Size)
		_, err = f.Read(buf)
		if err != nil {
			continue
		}

		files = append(files, buf)
	}

	return media.UploadInput{
		Files: files,
	}
}

func (r uploadRequest) validate() error {
	if len(r.Files) == 0 {
		return errWrongBody
	}
	return nil
}

type mediaItems struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id,omitempty"`
	ShopID   string `json:"shop_id,omitempty"`
	FileName string `json:"file_name"`
	Folder   string `json:"folder"`
	Status   string `json:"status"`
}

func newListResponse(medias []models.Media) []mediaItems {
	items := make([]mediaItems, 0, len(medias))

	for _, media := range medias {
		var shopID string
		if media.ShopID != primitive.NilObjectID {
			shopID = media.ShopID.Hex()
		}

		items = append(items, mediaItems{
			ID:       media.ID.Hex(),
			UserID:   media.UserID.Hex(),
			ShopID:   shopID,
			FileName: media.FileName,
			Folder:   media.Folder,
			Status:   media.Status,
		})
	}
	return items
}
