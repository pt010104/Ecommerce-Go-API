package usecase

import (
	"context"
	"errors"
	"testing"

	cloudinary "github.com/cloudinary/cloudinary-go"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/producer"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/log"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/stretchr/testify/require"
)

type mockDeps struct {
	repo  *media.MockRepository
	prod  *producer.MockProducer
	cloud *cloudinary.Cloudinary
}

func initUseCase(t *testing.T) (media.UseCase, mockDeps) {
	t.Helper()

	l := log.InitializeTestZapLogger()

	repo := media.NewMockRepository(t)
	prod := producer.NewMockProducer(t)
	cloud := &cloudinary.Cloudinary{}

	return New(l, repo, prod, *cloud), mockDeps{
		repo:  repo,
		prod:  prod,
		cloud: cloud,
	}
}

func TestList(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoList struct {
		isCalled bool
		input    media.ListOption
		output   []models.Media
		err      error
	}

	ids := []string{
		"6654408a9b657b844db56a74",
		"6654408a9b657b844db56a75",
		"6654408a9b657b844db56a76",
	}

	tcs := map[string]struct {
		input    media.ListInput
		mockRepo mockRepoList
		wantRes  []models.Media
		wantErr  error
	}{
		"success": {
			input: media.ListInput{
				GetFilter: media.GetFilter{
					IDs:    ids,
					Status: models.MediaStatusUploaded,
				},
			},
			mockRepo: mockRepoList{
				isCalled: true,
				input: media.ListOption{
					GetFilter: media.GetFilter{
						IDs:    ids,
						Status: models.MediaStatusUploaded,
					},
				},
				output: []models.Media{
					{
						ID:     mongo.ObjectIDFromHexOrNil(ids[0]),
						Status: models.MediaStatusUploaded,
					},
					{
						ID:     mongo.ObjectIDFromHexOrNil(ids[1]),
						Status: models.MediaStatusUploaded,
					},
					{
						ID:     mongo.ObjectIDFromHexOrNil(ids[2]),
						Status: models.MediaStatusUploaded,
					},
				},
				err: nil,
			},
			wantRes: []models.Media{
				{
					ID:     mongo.ObjectIDFromHexOrNil(ids[0]),
					Status: models.MediaStatusUploaded,
				},
				{
					ID:     mongo.ObjectIDFromHexOrNil(ids[1]),
					Status: models.MediaStatusUploaded,
				},
				{
					ID:     mongo.ObjectIDFromHexOrNil(ids[2]),
					Status: models.MediaStatusUploaded,
				},
			},
			wantErr: nil,
		},
		"success_with_empty_repo_output": {
			input: media.ListInput{
				GetFilter: media.GetFilter{
					IDs:    ids,
					Status: models.MediaStatusUploaded,
				},
			},
			mockRepo: mockRepoList{
				isCalled: true,
				input: media.ListOption{
					GetFilter: media.GetFilter{
						IDs:    ids,
						Status: models.MediaStatusUploaded,
					},
				},
				output: []models.Media{},
				err:    nil,
			},
			wantRes: []models.Media{},
			wantErr: nil,
		},
		"fail with empty IDs": {
			input: media.ListInput{
				GetFilter: media.GetFilter{
					IDs:    []string{},
					Status: models.MediaStatusUploaded,
				},
			},
			mockRepo: mockRepoList{
				isCalled: false,
			},
			wantRes: nil,
			wantErr: media.ErrRequireField,
		},
		"fail with invalid status": {
			input: media.ListInput{
				GetFilter: media.GetFilter{
					IDs:    ids,
					Status: "invalid_status",
				},
			},
			mockRepo: mockRepoList{
				isCalled: false,
			},
			wantRes: nil,
			wantErr: media.ErrInvalidStatus,
		},
		"fail with repo error": {
			input: media.ListInput{
				GetFilter: media.GetFilter{
					IDs:    ids,
					Status: models.MediaStatusUploaded,
				},
			},
			mockRepo: mockRepoList{
				isCalled: true,
				input: media.ListOption{
					GetFilter: media.GetFilter{
						IDs:    ids,
						Status: models.MediaStatusUploaded,
					},
				},
				output: nil,
				err:    errors.New("internal server error"),
			},
			wantRes: nil,
			wantErr: errors.New("internal server error"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			uc, deps := initUseCase(t)

			if tc.mockRepo.isCalled {
				deps.repo.EXPECT().List(ctx, scope, tc.mockRepo.input).
					Return(tc.mockRepo.output, tc.mockRepo.err)
			}

			res, err := uc.List(ctx, scope, tc.input)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantRes, res)
			}
		})
	}
}
