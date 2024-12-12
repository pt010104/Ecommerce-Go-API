package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) Upload(ctx context.Context, sc models.Scope, opt media.UploadInput) error {
	if len(opt.Files) == 0 {
		return media.ErrRequireField
	}

	for _, file := range opt.Files {
		uploadOpt := media.UploadOption{
			FileName: uc.generateFilename(sc.UserID),
			Folder:   uc.determineFolder(sc.UserID, sc.ShopID),
		}

		m, err := uc.repo.Create(ctx, sc, uploadOpt)
		if err != nil {
			uc.l.Errorf(ctx, "media.usecase.Upload.Create: %v", err)
			return err
		}

		err = uc.publishUploadMediaMessage(ctx, sc, file, uploadOpt, m.ID.Hex())
		if err != nil {
			uc.l.Errorf(ctx, "media.usecase.Upload.publishMediaMessage: %v", err)
			return err
		}
	}

	return nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (models.Media, error) {
	uc.l.Infof(ctx, "media.usecase.Detail: %v", id)
	m, err := uc.repo.Detail(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.Detail.repo.Detail: %v", err)
		return models.Media{}, err
	}

	return m, nil
}
