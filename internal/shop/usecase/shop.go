package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, input shop.CreateInput) (models.Shop, error) {
	if input.City == "" || input.Name == "" || input.Street == "" {
		return models.Shop{}, shop.ErrInvalidInput
	}

	_, err := uc.repo.Detail(ctx, sc, "")
	if err == nil {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", shop.ErrShopExist)
		return models.Shop{}, shop.ErrShopExist
	}

	if !util.IsValidPhone(input.Phone) {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", shop.ErrInvalidPhone)
		return models.Shop{}, shop.ErrInvalidPhone
	}

	opt := shop.CreateOption{
		Name:     input.Name,
		Alias:    util.BuildAlias(input.Name),
		City:     input.City,
		Street:   input.Street,
		District: input.District,
		Phone:    input.Phone,
	}

	sh, err := uc.repo.Create(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", err)
		return models.Shop{}, err
	}

	return sh, nil
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, input shop.GetInput) (shop.GetOutput, error) {
	opt := shop.GetOption{
		GetShopsFilter: input.GetShopsFilter,
		PagQuery:       input.PagQuery,
	}
	s, pag, err := uc.repo.Get(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Get: %v", err)
		return shop.GetOutput{}, err
	}

	return shop.GetOutput{
		Shops: s,
		Pag:   pag,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error) {
	s, err := uc.repo.Detail(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Detail: %v", err)
		return models.Shop{}, err
	}

	return s, nil
}
func (uc implUsecase) Delete(ctx context.Context, sc models.Scope, id string) (models.Shop, error) {

	shop1, err := uc.repo.FindByid(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "shop.delete.FindById:", err)
		return models.Shop{}, err
	}
	if shop1.UserID.Hex() != sc.UserID {
		uc.l.Errorf(ctx, "shop.usecase.Delete.Repodele", err)
		return models.Shop{}, shop.ErrNoPermissionToDelete
	}
	res, err := uc.repo.Delete(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Delete.Repodele", err)
		return models.Shop{}, shop.ErrShopDoesNotExist
	}

	return res, nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, input shop.UpdateInput) (models.Shop, error) {
	shop1, err := uc.repo.FindByid(ctx, sc, input.ID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.Update.FindById:", err)
		return models.Shop{}, err
	}
	if shop1.UserID.Hex() != sc.UserID {
		uc.l.Errorf(ctx, "shop.usecase.Update.RepoUpdate", err)
		return models.Shop{}, shop.ErrNoPermissionToUpdate
	}
	updateData := bson.M{}
	if input.Name != nil {
		updateData["name"] = *input.Name
	}
	if input.Alias != nil {
		updateData["alias"] = *input.Alias
	}
	if input.City != nil {
		updateData["city"] = *input.City
	}
	if input.Street != nil {
		updateData["street"] = *input.Street
	}
	if input.District != nil {
		updateData["district"] = *input.District
	}
	if input.Phone != nil {
		updateData["phone"] = *input.Phone
	}
	if input.Followers != nil {
		updateData["followers"] = *input.Followers
	}
	if input.AvgRate != nil {
		updateData["avg_rate"] = *input.AvgRate
	}

	updateData["updated_at"] = time.Now()
	updateOption := shop.UpdateOption{
		ID:         input.ID,
		UpdateData: updateData,
	}
	err = uc.repo.Update(ctx, sc, updateOption)
	if err != nil {
		return models.Shop{}, err
	}
	updatedShop, err := uc.repo.FindByid(ctx, sc, input.ID)
	if err != nil {
		return models.Shop{}, err
	}

	return updatedShop, nil
}
