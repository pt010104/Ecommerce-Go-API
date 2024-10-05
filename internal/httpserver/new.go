package httpserver

import (
	"github.com/gin-gonic/gin"
	pkgLog "github.com/pt010104/api-golang/pkg/log"
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
}

type Config struct {
	Port         int
	JWTSecretKey string
	Mode         string
	Database     mongo.Database
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
	}

}
