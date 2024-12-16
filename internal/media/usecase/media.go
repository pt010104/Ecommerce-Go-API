package usecase

import (
	"context"
	"slices"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) Upload(ctx context.Context, sc models.Scope, opt media.UploadInput) ([]models.Media, error) {
	if len(opt.Files) == 0 {
		return nil, media.ErrRequireField
	}

	medias := make([]models.Media, 0, len(opt.Files))

	for _, file := range opt.Files {
		uploadOpt := media.UploadOption{
			FileName: uc.generateFilename(sc.UserID),
			Folder:   uc.determineFolder(sc.UserID, sc.ShopID),
		}

		m, err := uc.repo.Create(ctx, sc, uploadOpt)
		if err != nil {
			uc.l.Errorf(ctx, "media.usecase.Upload.Create: %v", err)
			return nil, err
		}

		medias = append(medias, m)

		err = uc.publishUploadMediaMessage(ctx, sc, file, uploadOpt, m.ID.Hex())
		if err != nil {
			uc.l.Errorf(ctx, "media.usecase.Upload.publishMediaMessage: %v", err)
			return nil, err
		}
	}

	return medias, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (models.Media, error) {
	m, err := uc.repo.Detail(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.Detail.repo.Detail: %v", err)
		return models.Media{}, err
	}

	return m, nil
}

func (uc implUsecase) List(ctx context.Context, sc models.Scope, input media.ListInput) ([]models.Media, error) {
	if len(input.IDs) == 0 {
		return []models.Media{}, nil
	}

	if input.Status != "" && !slices.Contains(validStatus, input.Status) {
		return nil, media.ErrInvalidStatus
	} else {
		input.Status = models.MediaStatusUploaded
	}

	medias, err := uc.repo.List(ctx, sc, media.ListOption{
		GetFilter: input.GetFilter,
	})
	if err != nil {
		uc.l.Errorf(ctx, "media.usecase.List.repo.List: %v", err)
		return nil, err
	}

	return medias, nil
}
