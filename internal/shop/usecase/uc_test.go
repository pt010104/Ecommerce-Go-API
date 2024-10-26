package usecase

import (
	"testing"

	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/log"
)

type mockDeps struct {
	repo *shop.MockRepository
}

func initUseCase(t *testing.T) (shop.UseCase, mockDeps) {
	t.Helper()

	l := log.InitializeTestZapLogger()

	repo := shop.NewMockRepository(t)

	return New(l, repo), mockDeps{
		repo: repo,
	}
}
