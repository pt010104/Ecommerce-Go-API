package httpserver

import (
	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
	pkgLog "github.com/pt010104/api-golang/pkg/log"
	"github.com/pt010104/api-golang/pkg/rabbitmq"
	"github.com/pt010104/api-golang/pkg/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

const productionMode = "production"

var ginMode = gin.DebugMode

type HTTPServer struct {
	gin          *gin.Engine
	l            pkgLog.Logger
	port         int
	database     mongo.Database
	jwtSecretKey string
	mode         string
	redis        redis.Client
	amqpConn     rabbitmq.Connection
	cloudinary   cloudinary.Cloudinary
}

type Config struct {
	Port         int
	JWTSecretKey string
	Mode         string
	Database     mongo.Database
	Redis        redis.Client
	AMQPConn     rabbitmq.Connection
	Cloudinary   cloudinary.Cloudinary
}

func New(l pkgLog.Logger, cfg Config) *HTTPServer {
	if cfg.Mode == productionMode {
		ginMode = gin.ReleaseMode
	}

	gin.SetMode(ginMode)

	return &HTTPServer{
		l:            l,
		gin:          gin.Default(),
		database:     cfg.Database,
		jwtSecretKey: cfg.JWTSecretKey,
		port:         cfg.Port,
		mode:         cfg.Mode,
		redis:        cfg.Redis,
		amqpConn:     cfg.AMQPConn,
		cloudinary:   cfg.Cloudinary,
	}

}
