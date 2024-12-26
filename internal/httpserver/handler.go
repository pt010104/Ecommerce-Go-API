package httpserver

import (
	adminHTTP "github.com/pt010104/api-golang/internal/admins/delivery/http"
	cartHTTP "github.com/pt010104/api-golang/internal/cart/delivery/http"
	mediaHTTP "github.com/pt010104/api-golang/internal/media/delivery/http"
	"github.com/pt010104/api-golang/internal/middleware"
	orderHTTP "github.com/pt010104/api-golang/internal/order/delivery/http"
	shopHTTP "github.com/pt010104/api-golang/internal/shop/delivery/http"
	userHTTP "github.com/pt010104/api-golang/internal/user/delivery/http"
	voucherHTTP "github.com/pt010104/api-golang/internal/vouchers/delivery/http"

	adminRepo "github.com/pt010104/api-golang/internal/admins/repository/mongo"
	cartRepo "github.com/pt010104/api-golang/internal/cart/repository/mongo"
	cartUC "github.com/pt010104/api-golang/internal/cart/usecase"
	mediaRepo "github.com/pt010104/api-golang/internal/media/repository/mongo"
	orderRepo "github.com/pt010104/api-golang/internal/order/repository/mongo"
	redisOrderRepo "github.com/pt010104/api-golang/internal/order/repository/redis"
	orderUC "github.com/pt010104/api-golang/internal/order/usecase"
	shopRepo "github.com/pt010104/api-golang/internal/shop/repository/mongo"
	userRepo "github.com/pt010104/api-golang/internal/user/repository/mongo"
	redisUserRepo "github.com/pt010104/api-golang/internal/user/repository/redis"
	voucherRepo "github.com/pt010104/api-golang/internal/vouchers/repository/mongo"

	adminUC "github.com/pt010104/api-golang/internal/admins/usecase"
	emailUC "github.com/pt010104/api-golang/internal/email/usecase"
	mediaUC "github.com/pt010104/api-golang/internal/media/usecase"
	shopUC "github.com/pt010104/api-golang/internal/shop/usecase"
	userUC "github.com/pt010104/api-golang/internal/user/usecase"
	voucherUC "github.com/pt010104/api-golang/internal/vouchers/usecase"

	mediaProd "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq/producer"
	orderProd "github.com/pt010104/api-golang/internal/order/delivery/rabbitmq/producer"
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
	redisUserRepo := redisUserRepo.New(srv.l, srv.redis)
	redisOrderRepo := redisOrderRepo.New(srv.l, srv.redis)
	cartRepo := cartRepo.New(srv.l, srv.database)
	mediaRepo := mediaRepo.New(srv.l, srv.database)
	orderRepo := orderRepo.New(srv.l, srv.database)

	//Producer
	mediaProd := mediaProd.New(srv.l, srv.amqpConn)
	if err := mediaProd.Run(); err != nil {
		return err
	}
	orderProducer := orderProd.New(srv.l, srv.amqpConn)
	if err := orderProducer.Run(); err != nil {
		return err
	}
	//Usecase
	emailUC := emailUC.New(srv.l)
	mediaUC := mediaUC.New(srv.l, mediaRepo, mediaProd, srv.cloudinary)
	userUC := userUC.New(srv.l, userRepo, emailUC, redisUserRepo, mediaUC)
	shopUC := shopUC.New(srv.l, shopRepo, nil, userUC, mediaUC)
	adminUC := adminUC.New(adminRepo, srv.l, shopUC)
	shopUC.SetAdminUC(adminUC)
	cartUC := cartUC.New(srv.l, cartRepo, shopUC)
	voucherUC := voucherUC.New(voucherRepo, srv.l, shopUC)
	orderUC := orderUC.New(srv.l, orderRepo, shopUC, cartUC, redisOrderRepo, orderProducer, userUC)

	// Handlers
	userH := userHTTP.New(srv.l, userUC)
	shopH := shopHTTP.New(srv.l, shopUC)
	adminH := adminHTTP.New(srv.l, adminUC)
	voucherH := voucherHTTP.New(srv.l, voucherUC)
	cartH := cartHTTP.New(srv.l, cartUC)
	mediaH := mediaHTTP.New(srv.l, mediaUC)
	orderH := orderHTTP.New(srv.l, orderUC)
	mw := middleware.New(srv.l, shopUC, userUC)

	//Routes
	api := srv.gin.Group("/api/v1")

	userHTTP.MapRouters(api.Group("/users"), userH, mw)
	shopHTTP.MapRouters(api.Group("/shops"), shopH, mw)
	adminHTTP.MapRouters(api.Group("/admin"), adminH, mw)
	voucherHTTP.MapRouters(api.Group("/vouchers"), voucherH, mw)
	cartHTTP.MapRouters(api.Group("/carts"), cartH, mw)
	mediaHTTP.MapRouters(api.Group("/media"), mediaH, mw)
	orderHTTP.MapRouters(api.Group("/order"), orderH, mw)

	//Public routes
	shopHTTP.MapPublicRoutes(api.Group("/shops/products"), shopH)
	adminHTTP.MapPublicRoutes(api.Group("/categories"), adminH)

	return nil
}
