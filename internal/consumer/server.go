package consumer

import (
	"github.com/cloudinary/cloudinary-go"
	mediaConsumer "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/consumer"
	producer "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/producer"
	mediaRepo "github.com/pt010104/api-golang/internal/media/repository/mongo"
	mediaUC "github.com/pt010104/api-golang/internal/media/usecase"
	"github.com/pt010104/api-golang/pkg/log"
	"github.com/pt010104/api-golang/pkg/rabbitmq"
	"github.com/pt010104/api-golang/pkg/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

// Server is the consumer server
type Server struct {
	l          log.Logger
	conn       rabbitmq.Connection
	db         mongo.Database
	redis      redis.Client
	cloudinary cloudinary.Cloudinary
}

func NewServer(
	l log.Logger,
	conn rabbitmq.Connection,
	db mongo.Database,
	redis redis.Client,
	cloudinary cloudinary.Cloudinary,
) Server {
	return Server{
		l:          l,
		conn:       conn,
		db:         db,
		redis:      redis,
		cloudinary: cloudinary,
	}
}

func (s Server) Run() error {
	mediaRepo := mediaRepo.New(s.l, s.db)
	prod := producer.New(s.l, s.conn)
	if err := prod.Run(); err != nil {
		return err
	}
	mediaUC := mediaUC.New(s.l, mediaRepo, prod, s.cloudinary)

	// Consumers
	forever := make(chan bool)
	go func() {
		mediaConsumer.NewConsumer(s.l, &s.conn, mediaUC).Consume()
	}()

	select {
	case <-forever:
		panic("media consumer exited unexpectedly")
	}
}
