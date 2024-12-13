package usecase

import (
	"testing"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type mockDeps struct {
	repo    *shop.MockRepository
	adminUC *admins.MockUseCase
	userUC  *user.MockUseCase
	mediaUC *media.MockUseCase
}

func initUseCase(t *testing.T) (shop.UseCase, mockDeps) {
	t.Helper()

	l := log.InitializeTestZapLogger()

	repo := shop.NewMockRepository(t)
	adminUC := admins.NewMockUseCase(t)
	userUC := user.NewMockUseCase(t)
	mediaUC := media.NewMockUseCase(t)

	return New(l, repo, adminUC, userUC, mediaUC), mockDeps{
		repo:    repo,
		adminUC: adminUC,
		userUC:  userUC,
		mediaUC: mediaUC,
	}
}
