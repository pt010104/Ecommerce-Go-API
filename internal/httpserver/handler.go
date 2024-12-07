package httpserver

import (
	adminHTTP "github.com/pt010104/api-golang/internal/admins/delivery/http"
	cartHTTP "github.com/pt010104/api-golang/internal/cart/delivery/http"
	"github.com/pt010104/api-golang/internal/middleware"
	shopHTTP "github.com/pt010104/api-golang/internal/shop/delivery/http"
	userHTTP "github.com/pt010104/api-golang/internal/user/delivery/http"
	voucherHTTP "github.com/pt010104/api-golang/internal/vouchers/delivery/http"

	adminRepo "github.com/pt010104/api-golang/internal/admins/repository/mongo"
	cartRepo "github.com/pt010104/api-golang/internal/cart/repository/mongo"
	cartUC "github.com/pt010104/api-golang/internal/cart/usecase"
	shopRepo "github.com/pt010104/api-golang/internal/shop/repository/mongo"
	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"
	"github.com/pt010104/api-golang/internal/user/repository/mongo/redis"
	voucherRepo "github.com/pt010104/api-golang/internal/vouchers/repository/mongo"

	adminUC "github.com/pt010104/api-golang/internal/admins/usecase"
	emailUC "github.com/pt010104/api-golang/internal/email/usecase"
	shopUC "github.com/pt010104/api-golang/internal/shop/usecase"
	userUC "github.com/pt010104/api-golang/internal/user/usecase"
	voucherUC "github.com/pt010104/api-golang/internal/vouchers/usecase"

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
	voucherRepo := voucherRepo.New(srv.l, srv.database)
	redisRepo := redis.New(srv.l, srv.redis)
	cartRepo := cartRepo.New(srv.l, srv.database)
	//Usecase
	emailUC := emailUC.New(srv.l)
	userUC := userUC.New(srv.l, userRepo, emailUC, redisRepo)
	shopUC := shopUC.New(srv.l, shopRepo, nil)
	adminUC := adminUC.New(adminRepo, srv.l, shopUC)
	shopUC.SetAdminUC(adminUC)
	cartUC := cartUC.New(srv.l, cartRepo, shopUC)
	voucherUC := voucherUC.New(voucherRepo, srv.l, shopUC)

	// Handlers
	userH := userHTTP.New(srv.l, userUC)
	shopH := shopHTTP.New(srv.l, shopUC)
	adminH := adminHTTP.New(srv.l, adminUC)
	voucherH := voucherHTTP.New(srv.l, voucherUC)
	cartH := cartHTTP.New(srv.l, cartUC)
	mw := middleware.New(srv.l, userRepo, shopUC, userUC)

	//Routes
	api := srv.gin.Group("/api/v1")

	userHTTP.MapRouters(api.Group("/users"), userH, mw)
	shopHTTP.MapRouters(api.Group("/shops"), shopH, mw)
	adminHTTP.MapRouters(api.Group("/admin"), adminH, mw)
	voucherHTTP.MapRouters(api.Group("/vouchers"), voucherH, mw)
	cartHTTP.MapRouters(api.Group("/carts"), cartH, mw)

	return nil
}
