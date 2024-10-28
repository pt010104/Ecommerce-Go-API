package usecase

import (
	"context"
	"testing"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateInventory(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoCreate struct {
		isCalled bool
		input    shop.CreateInventoryOption
		output   models.Inventory
		err      error
	}

	id := primitive.NewObjectID()
	reorderLevel := uint(5)
	reorderLevelPtr := &reorderLevel

	tcs := map[string]struct {
		input    shop.CreateInventoryInput
		mockRepo mockRepoCreate
		wantRes  models.Inventory
		wantErr  error
	}{
		"success without ReorderLevel": {
			input: shop.CreateInventoryInput{
				ProductID:  id,
				StockLevel: 10,
			},
			mockRepo: mockRepoCreate{
				isCalled: true,
				input: shop.CreateInventoryOption{
					ProductID:  id,
					StockLevel: 10,
				},
				output: models.Inventory{
					ProductID:  id,
					StockLevel: 10,
				},
				err: nil,
			},
			wantRes: models.Inventory{
				ProductID:  id,
				StockLevel: 10,
			},

			wantErr: nil,
		},
		"success with ReorderLevel but without ReorderQuantity": {
			input: shop.CreateInventoryInput{
				ProductID:    id,
				StockLevel:   10,
				ReorderLevel: reorderLevelPtr,
			},
			mockRepo: mockRepoCreate{
				isCalled: true,
				input: shop.CreateInventoryOption{
					ProductID:       id,
					StockLevel:      10,
					ReorderLevel:    reorderLevelPtr,
					ReorderQuantity: nil,
				},
				output: models.Inventory{
					ProductID:       id,
					StockLevel:      10,
					ReorderLevel:    nil,
					ReorderQuantity: nil,
				},
				err: nil,
			},
			wantRes: models.Inventory{
				ProductID:       id,
				StockLevel:      10,
				ReorderLevel:    nil,
				ReorderQuantity: nil,
			},
			wantErr: nil,
		},
		"success with ReorderLevel with ReorderQuantity": {
			input: shop.CreateInventoryInput{
				ProductID:       id,
				StockLevel:      10,
				ReorderLevel:    reorderLevelPtr,
				ReorderQuantity: reorderLevelPtr,
			},
			mockRepo: mockRepoCreate{
				isCalled: true,
				input: shop.CreateInventoryOption{
					ProductID:       id,
					StockLevel:      10,
					ReorderLevel:    reorderLevelPtr,
					ReorderQuantity: reorderLevelPtr,
				},
				output: models.Inventory{
					ProductID:       id,
					StockLevel:      10,
					ReorderLevel:    reorderLevelPtr,
					ReorderQuantity: reorderLevelPtr,
				},
				err: nil,
			},
			wantRes: models.Inventory{
				ProductID:       id,
				StockLevel:      10,
				ReorderLevel:    reorderLevelPtr,
				ReorderQuantity: reorderLevelPtr,
			},
			wantErr: nil,
		},
		"fail with not exist productID": {
			input: shop.CreateInventoryInput{
				ProductID:  id,
				StockLevel: 10,
			},
			mockRepo: mockRepoCreate{
				isCalled: true,
				input: shop.CreateInventoryOption{
					ProductID:  id,
					StockLevel: 10,
				},
				output: models.Inventory{},
				err:    mongo.ErrNoDocuments,
			},
			wantRes: models.Inventory{},
			wantErr: mongo.ErrNoDocuments,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			uc, deps := initUseCase(t)

			if tc.mockRepo.isCalled {
				deps.repo.EXPECT().CreateInventory(ctx, scope, tc.mockRepo.input).
					Return(
						tc.mockRepo.output,
						tc.mockRepo.err,
					)
			}

			res, err := uc.CreateInventory(ctx, scope, tc.input)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantRes, res)
			}

		})
	}

}

func TestDetailInventory(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoDetail struct {
		isCalled  bool
		productID primitive.ObjectID
		output    models.Inventory
		err       error
	}

	id := primitive.NewObjectID()

	tcs := map[string]struct {
		productID primitive.ObjectID
		mockRepo  mockRepoDetail
		wantRes   models.Inventory
		wantErr   error
	}{
		"success": {
			productID: id,
			mockRepo: mockRepoDetail{
				isCalled:  true,
				productID: id,
				output: models.Inventory{
					ProductID: id,
				},
				err: nil,
			},
			wantRes: models.Inventory{
				ProductID: id,
			},
			wantErr: nil,
		},
		"fail with not exist productID": {
			productID: id,
			mockRepo: mockRepoDetail{
				isCalled:  true,
				productID: id,
				output:    models.Inventory{},
				err:       mongo.ErrNoDocuments,
			},
			wantRes: models.Inventory{},
			wantErr: mongo.ErrNoDocuments,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			uc, deps := initUseCase(t)

			if tc.mockRepo.isCalled {
				deps.repo.EXPECT().DetailInventory(ctx, scope, tc.mockRepo.productID).
					Return(
						tc.mockRepo.output,
						tc.mockRepo.err,
					)
			}

			res, err := uc.DetailInventory(ctx, scope, tc.productID)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantRes, res)
			}

		})
	}
}

func TestUpdateInventory(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoUpdate struct {
		isCalled bool
		input    shop.UpdateInventoryOption
		output   models.Inventory
		err      error
	}

	type mockRepoDetail struct {
		isCalled  bool
		productID primitive.ObjectID
		output    models.Inventory
		err       error
	}

	type mockRepo struct {
		mockRepoDetail
		mockRepoUpdate
	}

	id := primitive.NewObjectID()
	tmp := uint(5)
	uintPtr := &tmp

	tcs := map[string]struct {
		input    shop.UpdateInventoryInput
		mockRepo mockRepo
		wantRes  models.Inventory
		wantErr  error
	}{
		"success without ReorderLevel": {
			input: shop.UpdateInventoryInput{
				ProductID:  id,
				StockLevel: uintPtr,
			},
			mockRepo: mockRepo{
				mockRepoDetail: mockRepoDetail{
					isCalled:  true,
					productID: id,
					output: models.Inventory{
						ProductID:  id,
						StockLevel: 10,
					},
					err: nil,
				},
				mockRepoUpdate: mockRepoUpdate{
					isCalled: true,
					input: shop.UpdateInventoryOption{
						Model: models.Inventory{
							ProductID:  id,
							StockLevel: 10,
						},
						StockLevel: uintPtr,
					},
					output: models.Inventory{
						ProductID:  id,
						StockLevel: 5,
					},
					err: nil,
				},
			},
			wantRes: models.Inventory{
				ProductID:  id,
				StockLevel: 5,
			},
			wantErr: nil,
		},
		"success with ReorderLevel but without ReorderQuantity": {
			input: shop.UpdateInventoryInput{
				ProductID:    id,
				StockLevel:   uintPtr,
				ReorderLevel: uintPtr,
			},
			mockRepo: mockRepo{
				mockRepoDetail: mockRepoDetail{
					isCalled:  true,
					productID: id,
					output: models.Inventory{
						ProductID:  id,
						StockLevel: 10,
					},
					err: nil,
				},
				mockRepoUpdate: mockRepoUpdate{
					isCalled: true,
					input: shop.UpdateInventoryOption{
						Model: models.Inventory{
							ProductID:  id,
							StockLevel: 10,
						},
						StockLevel:   uintPtr,
						ReorderLevel: uintPtr,
					},
					output: models.Inventory{
						ProductID:    id,
						StockLevel:   5,
						ReorderLevel: uintPtr,
					},
					err: nil,
				},
			},
			wantRes: models.Inventory{
				ProductID:    id,
				StockLevel:   5,
				ReorderLevel: uintPtr,
			},
			wantErr: nil,
		},
		"success with ReorderLevel and ReorderQuantity": {
			input: shop.UpdateInventoryInput{
				ProductID:       id,
				StockLevel:      uintPtr,
				ReorderLevel:    uintPtr,
				ReorderQuantity: uintPtr,
			},
			mockRepo: mockRepo{
				mockRepoDetail: mockRepoDetail{
					isCalled:  true,
					productID: id,
					output: models.Inventory{
						ProductID:  id,
						StockLevel: 10,
					},
					err: nil,
				},
				mockRepoUpdate: mockRepoUpdate{
					isCalled: true,
					input: shop.UpdateInventoryOption{
						Model: models.Inventory{
							ProductID:  id,
							StockLevel: 10,
						},
						StockLevel:      uintPtr,
						ReorderLevel:    uintPtr,
						ReorderQuantity: uintPtr,
					},
					output: models.Inventory{
						ProductID:       id,
						StockLevel:      5,
						ReorderLevel:    uintPtr,
						ReorderQuantity: uintPtr,
					},
					err: nil,
				},
			},
			wantRes: models.Inventory{
				ProductID:       id,
				StockLevel:      5,
				ReorderLevel:    uintPtr,
				ReorderQuantity: uintPtr,
			},
			wantErr: nil,
		},
		"fail with not exist productID": {
			input: shop.UpdateInventoryInput{
				ProductID: id,
			},
			mockRepo: mockRepo{
				mockRepoDetail: mockRepoDetail{
					isCalled:  true,
					productID: id,
					output:    models.Inventory{},
					err:       mongo.ErrNoDocuments,
				},
			},
			wantRes: models.Inventory{},
			wantErr: mongo.ErrNoDocuments,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			uc, deps := initUseCase(t)

			if tc.mockRepo.mockRepoDetail.isCalled {
				deps.repo.EXPECT().DetailInventory(ctx, scope, tc.mockRepo.productID).
					Return(
						tc.mockRepo.mockRepoDetail.output,
						tc.mockRepo.mockRepoDetail.err,
					)
			}

			if tc.mockRepo.mockRepoUpdate.isCalled {
				deps.repo.EXPECT().UpdateInventory(ctx, scope, tc.mockRepo.mockRepoUpdate.input).
					Return(
						tc.mockRepo.mockRepoUpdate.output,
						tc.mockRepo.mockRepoUpdate.err,
					)
			}

			res, err := uc.UpdateInventory(ctx, scope, tc.input)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.mockRepo.mockRepoUpdate.output, res)
			}

		})
	}
}

func TestListlInventory(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoList struct {
		isCalled   bool
		productIDs []primitive.ObjectID
		output     []models.Inventory
		err        error
	}

	ids := []primitive.ObjectID{
		mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
		mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
		mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a76"),
	}

	tcs := map[string]struct {
		productIDs []primitive.ObjectID
		mockRepo   mockRepoList
		wantRes    []models.Inventory
		wantErr    error
	}{
		"success": {
			productIDs: ids,
			mockRepo: mockRepoList{
				isCalled:   true,
				productIDs: ids,
				output: []models.Inventory{
					{
						ProductID: ids[0],
					},
					{
						ProductID: ids[1],
					},
					{
						ProductID: ids[2],
					},
				},
				err: nil,
			},
			wantRes: []models.Inventory{
				{
					ProductID: ids[0],
				},
				{
					ProductID: ids[1],
				},
				{
					ProductID: ids[2],
				},
			},
			wantErr: nil,
		},
		"success with not exist 1 productID": {
			productIDs: ids,
			mockRepo: mockRepoList{
				isCalled:   true,
				productIDs: ids,
				output: []models.Inventory{
					{
						ProductID: ids[0],
					},
					{
						ProductID: ids[1],
					},
				},
				err: nil,
			},
			wantRes: []models.Inventory{
				{
					ProductID: ids[0],
				},
				{
					ProductID: ids[1],
				},
			},
			wantErr: mongo.ErrNoDocuments,
		},
		"success with not exist any productID": {
			productIDs: ids,
			mockRepo: mockRepoList{
				isCalled:   true,
				productIDs: ids,
				output:     []models.Inventory{},
				err:        nil,
			},
			wantRes: []models.Inventory{},
			wantErr: mongo.ErrNoDocuments,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			uc, deps := initUseCase(t)

			if tc.mockRepo.isCalled {
				deps.repo.EXPECT().ListInventory(ctx, scope, tc.mockRepo.productIDs).
					Return(
						tc.mockRepo.output,
						tc.mockRepo.err,
					)
			}

			res, err := uc.ListInventory(ctx, scope, tc.productIDs)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantRes, res)
			}

		})
	}
}

func TestDeleteInventory(t *testing.T) {
	scope := models.Scope{
		UserID:    "test",
		SessionID: "test",
		Role:      0,
	}

	type mockRepoDelete struct {
		isCalled   bool
		productIDs []primitive.ObjectID
		err        error
	}

	ids := []primitive.ObjectID{
		mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a74"),
		mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a75"),
		mongo.ObjectIDFromHexOrNil("6654408a9b657b844db56a76"),
	}

	tcs := map[string]struct {
		productIDs []primitive.ObjectID
		mockRepo   mockRepoDelete
		wantErr    error
	}{
		"success": {
			productIDs: ids,
			mockRepo: mockRepoDelete{
				isCalled:   true,
				productIDs: ids,
				err:        nil,
			},
			wantErr: nil,
		},
		"success with not exist 1 productID": {
			productIDs: ids,
			mockRepo: mockRepoDelete{
				isCalled:   true,
				productIDs: ids,
				err:        nil,
			},
			wantErr: nil,
		},
		"success with not exist any productID": {
			productIDs: ids,
			mockRepo: mockRepoDelete{
				isCalled:   true,
				productIDs: ids,
				err:        nil,
			},
			wantErr: nil,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			uc, deps := initUseCase(t)

			if tc.mockRepo.isCalled {
				deps.repo.EXPECT().DeleteInventory(ctx, scope, tc.mockRepo.productIDs).
					Return(
						tc.mockRepo.err,
					)
			}

			err := uc.DeleteInventory(ctx, scope, tc.productIDs)
			if err != nil {
				require.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
			}

		})
	}
}
