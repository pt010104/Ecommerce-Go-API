package httpserver

import (
	"github.com/pt010104/api-golang/internal/middleware"
	userHTTP "github.com/pt010104/api-golang/internal/user/delivery/http"

	emailUC "github.com/pt010104/api-golang/internal/email/usecase"
	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"
	userUC "github.com/pt010104/api-golang/internal/user/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (srv HTTPServer) mapHandlers() error {

	//swagger
	srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Repo
	userRepo := userRepo.New(srv.l, srv.database)
	emailUC := emailUC.New(srv.l)
	//Usecase
	userUC := userUC.New(srv.l, userRepo, emailUC)

	// Handlers
	userH := userHTTP.New(srv.l, userUC)

	mw := middleware.New(srv.l, userRepo)

	//Routes
	api := srv.gin.Group("/api/v1")

	userHTTP.MapRouters(api.Group("/users"), userH, mw)

	return nil
}
