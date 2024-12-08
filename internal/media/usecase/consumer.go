package usecase

import (
	"context"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc *implUsecase) ConsumeUploadMessage(ctx context.Context, sc models.Scope, input media.ConsumeUploadMsgInput) error {
	status := models.MediaStatusFailed
	med, err := uc.repo.Detail(ctx, sc, input.ID)
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.repo.Detail: %v", err)
		return err
	}

	res, err := uc.cloud.Upload.Upload(ctx, input.File, uploader.UploadParams{
		Folder:   input.FolderName,
		PublicID: input.ID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.cloud.Upload.Upload: %v", err)
		status = models.MediaStatusFailed
	}

	_, err = uc.repo.Update(ctx, sc, input.ID, media.UpdateOption{
		Model:  med,
		URL:    res.URL,
		Status: status,
	})
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.ConsumeUploadMessage.repo.Update: %v", err)
		return err
	}

	return nil
}
