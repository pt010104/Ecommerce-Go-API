package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) Upload(ctx context.Context, sc models.Scope, opt media.UploadInput) error {
	if (opt.UserID.IsZero() && opt.ShopID.IsZero()) || len(opt.Files) == 0 {
		return media.ErrRequireField
	}

	for _, file := range opt.Files {
		uploadOpt := media.UploadOption{
			UserID:   opt.UserID,
			ShopID:   opt.ShopID,
			FileName: uc.generateFilename(opt.UserID),
			Folder:   uc.determineFolder(opt.UserID, opt.ShopID),
		}

		err := uc.repo.Create(ctx, sc, uploadOpt)
		if err != nil {
			uc.l.Errorf(ctx, "media.usecase.Upload.Create: %v", err)
			return err
		}

		err = uc.publishUploadMediaMessage(ctx, file, uploadOpt)
		if err != nil {
			uc.l.Errorf(ctx, "media.usecase.Upload.publishMediaMessage: %v", err)
			return err
		}
	}

	return nil
}
