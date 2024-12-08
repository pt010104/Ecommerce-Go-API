package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
)

func (uc implUsecase) publishUploadMediaMessage(ctx context.Context, file []byte, uploadOpt media.UploadOption) error {
	message := rabbitmq.UploadMessage{
		File:       file,
		UserID:     uploadOpt.UserID,
		ShopID:     uploadOpt.ShopID,
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
