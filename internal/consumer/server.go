package consumer

import (
	"github.com/cloudinary/cloudinary-go"
	adminRepo "github.com/pt010104/api-golang/internal/admins/repository/mongo"
	adminUC "github.com/pt010104/api-golang/internal/admins/usecase"
	cartRepo "github.com/pt010104/api-golang/internal/cart/repository/mongo"
	cartUC "github.com/pt010104/api-golang/internal/cart/usecase"
	emailUC "github.com/pt010104/api-golang/internal/email/usecase"
	mediaConsumer "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/consumer"
	producer "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/producer"
	mediaRepo "github.com/pt010104/api-golang/internal/media/repository/mongo"
	mediaUC "github.com/pt010104/api-golang/internal/media/usecase"
	orderConsumer "github.com/pt010104/api-golang/internal/order/delivery/rabbitmq/consumer"
	orderProd "github.com/pt010104/api-golang/internal/order/delivery/rabbitmq/producer"
	orderRepo "github.com/pt010104/api-golang/internal/order/repository/mongo"
	redisOrderRepo "github.com/pt010104/api-golang/internal/order/repository/redis"
	orderUC "github.com/pt010104/api-golang/internal/order/usecase"
	shopRepo "github.com/pt010104/api-golang/internal/shop/repository/mongo"
	shopUC "github.com/pt010104/api-golang/internal/shop/usecase"
	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"
	redisUserRepo "github.com/pt010104/api-golang/internal/user/repository/redis"
	userUC "github.com/pt010104/api-golang/internal/user/usecase"
	voucherRepo "github.com/pt010104/api-golang/internal/vouchers/repository/mongo"
	voucherUC "github.com/pt010104/api-golang/internal/vouchers/usecase"
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
	orderRepo := orderRepo.New(s.l, s.db)
	shopRepo := shopRepo.New(s.l, s.db)
	userRepo := userRepo.New(s.l, s.db)
	cartRepo := cartRepo.New(s.l, s.db)
	redisOrderRepo := redisOrderRepo.New(s.l, s.redis)
	redisUserRepo := redisUserRepo.New(s.l, s.redis)
	adminRepo := adminRepo.New(s.l, s.db)
	voucherRepo := voucherRepo.New(s.l, s.db)

	prod := producer.New(s.l, s.conn)
	if err := prod.Run(); err != nil {
		return err
	}

	orderProducer := orderProd.New(s.l, s.conn)
	if err := orderProducer.Run(); err != nil {
		return err
	}

	emailUC := emailUC.New(s.l)
	mediaUC := mediaUC.New(s.l, mediaRepo, prod, s.cloudinary)
	userUC := userUC.New(s.l, userRepo, emailUC, redisUserRepo, mediaUC)
	shopUC := shopUC.New(s.l, shopRepo, nil, userUC, mediaUC)
	cartUC := cartUC.New(s.l, cartRepo, shopUC)
	voucherUC := voucherUC.New(voucherRepo, s.l, shopUC)
	orderUC := orderUC.New(s.l, orderRepo, shopUC, cartUC, redisOrderRepo, orderProducer, emailUC, userUC, voucherUC)
	adminUC := adminUC.New(adminRepo, s.l, shopUC)

	shopUC.SetAdminUC(adminUC)

	// Consumers
	forever := make(chan bool)
	go func() {
		mediaConsumer.NewConsumer(s.l, &s.conn, mediaUC).Consume()
		orderConsumer.NewConsumer(s.l, &s.conn, orderUC).Consume()
	}()

	select {
	case <-forever:
		panic("consumer exited unexpectedly")
	}
}
