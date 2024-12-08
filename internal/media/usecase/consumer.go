package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) getMimeType(file []byte) string {
	return http.DetectContentType(file)
}

func (uc implUsecase) ConsumeUploadMsg(ctx context.Context, sc models.Scope, input media.ConsumeUploadMsgInput) error {
	med, err := uc.repo.Detail(ctx, sc, input.ID)
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.repo.Detail: %v", err)
		return err
	}

	base64Data := base64.StdEncoding.EncodeToString(input.File)
	dataURI := fmt.Sprintf("data:%s;base64,%s", uc.getMimeType(input.File), base64Data)

	res, err := uc.cloud.Upload.Upload(ctx, dataURI, uploader.UploadParams{
		Folder:   input.FolderName,
		PublicID: input.FileName,
	})
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.cloud.Upload.Upload: %v", err)
		_, updateErr := uc.repo.Update(ctx, sc, input.ID, media.UpdateOption{
			Model:  med,
			Status: models.MediaStatusFailed,
		})
		if updateErr != nil {
			uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.repo.Update: %v", updateErr)
			return updateErr
		}
		return err
	}

	_, err = uc.repo.Update(ctx, sc, input.ID, media.UpdateOption{
		Model:  med,
		URL:    res.URL,
		Status: models.MediaStatusUploaded,
	})
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.repo.Update: %v", err)
		return err
	}

	return nil
}
