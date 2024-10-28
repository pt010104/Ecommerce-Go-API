package usecase

import (
	"context"
	"testing"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/stretchr/testify/require"
)

func TestUpdateShop(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoUpdate struct {
		isCalled bool
		input    []shop.UpdateOption
		output   []models.Shop
		err      error
	}

	type mockRepoList struct {
		isCalled bool
		intput   shop.GetOption
		output   []models.Shop
		err      error
	}

	type mockRepo struct {
		mockRepoUpdate
		mockRepoList
	}

	ids := []string{
		"6654408a9b657b844db56a74",
		"6654408a9b657b844db56a75",
		"6654408a9b657b844db56a76",
	}

	updateOptions := []shop.UpdateOption{
		{
			Model: models.Shop{
				ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
				IsVerified: false,
			},
			IsVerified: true,
		},
		{
			Model: models.Shop{
				ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
				IsVerified: false,
			},
			IsVerified: true,
		},
		{
			Model: models.Shop{
				ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a76"),
				IsVerified: false,
			},
			IsVerified: true,
		},
	}

	tcs := map[string]struct {
		input    shop.UpdateInput
		mockRepo mockRepo
		wantRes  []models.Shop
		wantErr  error
	}{
		"success with many ids": {
			input: shop.UpdateInput{
				ShopIDs:    ids,
				IsVerified: true,
			},
			mockRepo: mockRepo{
				mockRepoUpdate: mockRepoUpdate{
					isCalled: true,
					input:    updateOptions,
					output: []models.Shop{
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
							IsVerified: true,
						},
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
							IsVerified: true,
						},
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a76"),
							IsVerified: true,
						},
					},
				},
				mockRepoList: mockRepoList{
					isCalled: true,
					intput: shop.GetOption{
						GetShopsFilter: shop.GetShopsFilter{
							IDs: ids,
						},
					},

					output: []models.Shop{
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
							IsVerified: false,
						},
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
							IsVerified: false,
						},
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a76"),
							IsVerified: false,
						},
					},
				},
			},
			wantRes: []models.Shop{
				{
					ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
					IsVerified: true,
				},
				{
					ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
					IsVerified: true,
				},
				{
					ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a76"),
					IsVerified: true,
				},
			},
			wantErr: nil,
		},
		"success with 1 id": {
			input: shop.UpdateInput{
				ShopIDs:    []string{ids[0]},
				IsVerified: true,
			},
			mockRepo: mockRepo{
				mockRepoUpdate: mockRepoUpdate{
					isCalled: true,
					input:    []shop.UpdateOption{updateOptions[0]},
					output: []models.Shop{
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
							IsVerified: true,
						},
					},
				},
				mockRepoList: mockRepoList{
					isCalled: true,
					intput: shop.GetOption{
						GetShopsFilter: shop.GetShopsFilter{
							IDs: []string{ids[0]},
						},
					},

					output: []models.Shop{
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
							IsVerified: false,
						},
					},
				},
			},
			wantRes: []models.Shop{
				{
					ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
					IsVerified: true,
				},
			},
			wantErr: nil,
		},
		"success with lack of 1 id": {
			input: shop.UpdateInput{
				ShopIDs:    ids[:len(ids)-1],
				IsVerified: true,
			},
			mockRepo: mockRepo{
				mockRepoUpdate: mockRepoUpdate{
					isCalled: true,
					input:    updateOptions[:2],
					output: []models.Shop{
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
							IsVerified: true,
						},
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
							IsVerified: true,
						},
					},
				},
				mockRepoList: mockRepoList{
					isCalled: true,
					intput: shop.GetOption{
						GetShopsFilter: shop.GetShopsFilter{
							IDs: ids[:len(ids)-1],
						},
					},
					output: []models.Shop{
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
							IsVerified: false,
						},
						{
							ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
							IsVerified: false,
						},
					},
				},
			},
			wantRes: []models.Shop{
				{
					ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
					IsVerified: true,
				},
				{
					ID:         mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
					IsVerified: true,
				},
			},
			wantErr: nil,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			uc, deps := initUseCase(t)

			if tc.mockRepo.mockRepoList.isCalled {
				deps.repo.EXPECT().ListShop(ctx, scope, tc.mockRepo.mockRepoList.intput).
					Return(
						tc.mockRepo.mockRepoList.output,
						tc.mockRepo.mockRepoList.err,
					)
			}

			if tc.mockRepo.mockRepoUpdate.isCalled {
				for i := range tc.mockRepo.mockRepoUpdate.input {
					deps.repo.EXPECT().UpdateShop(ctx, scope, tc.mockRepo.mockRepoUpdate.input[i]).
						Return(
							tc.mockRepo.mockRepoUpdate.output[i],
							tc.mockRepo.mockRepoUpdate.err,
						)
				}
			}

			res, err := uc.Update(ctx, scope, tc.input)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantRes, res)
			}

		})
	}

}
