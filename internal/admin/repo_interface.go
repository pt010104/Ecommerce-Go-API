package admin

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repo interface {
	CreateCategory(ctx context.Context, sc models.Scope, opt CreateCategoryOption) (models.Category, error)
}
