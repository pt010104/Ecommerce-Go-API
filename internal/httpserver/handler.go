package httpserver

import (
	"github.com/pt010104/api-golang/internal/middleware"
	shopHTTP "github.com/pt010104/api-golang/internal/shop/delivery/http"
	userHTTP "github.com/pt010104/api-golang/internal/user/delivery/http"

	shopRepo "github.com/pt010104/api-golang/internal/shop/repository/mongo"
	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"

	emailUC "github.com/pt010104/api-golang/internal/email/usecase"
	shopUC "github.com/pt010104/api-golang/internal/shop/usecase"
	userUC "github.com/pt010104/api-golang/internal/user/usecase"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/pt010104/api-golang/docs"
)

func (srv HTTPServer) mapHandlers() error {

	//swagger
	srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Repo
	userRepo := userRepo.New(srv.l, srv.database)
	shopRepo := shopRepo.New(srv.l, srv.database)

	//Usecase
	emailUC := emailUC.New(srv.l)
	userUC := userUC.New(srv.l, userRepo, emailUC)
	shopUC := shopUC.New(srv.l, shopRepo)

	// Handlers
	userH := userHTTP.New(srv.l, userUC)
	shopH := shopHTTP.New(srv.l, shopUC)

	mw := middleware.New(srv.l, userRepo)

	//Routes
	api := srv.gin.Group("/api/v1")

	userHTTP.MapRouters(api.Group("/users"), userH, mw)
	shopHTTP.MapRouters(api.Group("/shops"), shopH, mw)

	return nil
}
