package usecase

import (
	"github.com/pt010104/api-golang/pkg/log"

	cloudinary "github.com/cloudinary/cloudinary-go"
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/producer"
)

type implUsecase struct {
	l     log.Logger
	repo  media.Repository
	prod  producer.Producer
	cloud cloudinary.Cloudinary
}

func New(l log.Logger, repo media.Repository, prod producer.Producer, cloud cloudinary.Cloudinary) media.UseCase {
	return &implUsecase{
		l:    l,
		repo: repo,
		prod: prod,
	}
}
