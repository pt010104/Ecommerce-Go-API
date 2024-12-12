package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) publishUploadMediaMessage(ctx context.Context, sc models.Scope, file []byte, uploadOpt media.UploadOption, id string) error {
	message := rabbitmq.UploadMessage{
		File:       file,
		ID:         id,
		UserID:     sc.UserID,
		ShopID:     sc.ShopID,
		FolderName: uploadOpt.Folder,
		FileName:   uploadOpt.FileName,
	}

	err := uc.prod.PublishUploadMsg(ctx, message)
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.publishUploadMediaMessage.prod.PublishUploadMsg: %v", err)
		return err
	}

	return nil
}
