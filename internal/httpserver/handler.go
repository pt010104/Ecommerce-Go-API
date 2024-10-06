package httpserver

import (
	"github.com/pt010104/api-golang/internal/middleware"
	userHTTP "github.com/pt010104/api-golang/internal/user/delivery/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	userUC "github.com/pt010104/api-golang/internal/user/usecase"

	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"
)

func (srv HTTPServer) mapHandlers() error {

	//swagger
	srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Repo
	userRepo := userRepo.New(srv.l, srv.database)

	//Usecase
	userUC := userUC.New(srv.l, userRepo)

	// Handlers
	userH := userHTTP.New(srv.l, userUC)

	mw := middleware.New(srv.l, userRepo)

	//Routes
	api := srv.gin.Group("/api/v1")

	userHTTP.MapRouters(api.Group("/users"), userH, mw)

	return nil
}
