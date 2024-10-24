package httpserver

import (
	adminHTTP "github.com/pt010104/api-golang/internal/admin/delivery/http"
	"github.com/pt010104/api-golang/internal/middleware"
	shopHTTP "github.com/pt010104/api-golang/internal/shop/delivery/http"
	userHTTP "github.com/pt010104/api-golang/internal/user/delivery/http"

	adminRepo "github.com/pt010104/api-golang/internal/admin/repository/mongo"
	shopRepo "github.com/pt010104/api-golang/internal/shop/repository/mongo"
	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"

	adminUC "github.com/pt010104/api-golang/internal/admin/usecase"
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
	adminRepo := adminRepo.New(srv.l, srv.database)

	//Usecase
	emailUC := emailUC.New(srv.l)
	userUC := userUC.New(srv.l, userRepo, emailUC)
	shopUC := shopUC.New(srv.l, shopRepo)
	adminUC := adminUC.New(adminRepo, srv.l)

	// Handlers
	userH := userHTTP.New(srv.l, userUC)
	shopH := shopHTTP.New(srv.l, shopUC)
	adminH := adminHTTP.New(srv.l, adminUC)

	mw := middleware.New(srv.l, userRepo)

	//Routes
	api := srv.gin.Group("/api/v1")

	userHTTP.MapRouters(api.Group("/users"), userH, mw)
	shopHTTP.MapRouters(api.Group("/shops"), shopH, mw)
	adminHTTP.MapRouters(api.Group("/admin"), adminH, mw)

	return nil
}
