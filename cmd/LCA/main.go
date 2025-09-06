package main

import (
	"github.com/wang900115/LCA/internal/adapter/controller"
	response "github.com/wang900115/LCA/internal/adapter/controller/response/json"
	"github.com/wang900115/LCA/internal/adapter/middleware"
	corsMid "github.com/wang900115/LCA/internal/adapter/middleware/cors"
	jwtMid "github.com/wang900115/LCA/internal/adapter/middleware/jwt"
	redisrate "github.com/wang900115/LCA/internal/adapter/middleware/redis_rate"
	"github.com/wang900115/LCA/internal/bootstrap"
	gormimplement "github.com/wang900115/LCA/internal/implement/gorm"
	redisimplement "github.com/wang900115/LCA/internal/implement/redis"

	// redisrate "LCA/internal/adapter/gin/middleware/redis_rate"
	secureheader "github.com/wang900115/LCA/internal/adapter/middleware/secure_header"
	"github.com/wang900115/LCA/internal/adapter/router"
	"github.com/wang900115/LCA/internal/adapter/websocket/connection"
	"github.com/wang900115/LCA/internal/application/usecase"
)

func main() {
	conf := bootstrap.NewConfig()

	redispool := bootstrap.NewRedisPool(bootstrap.NewRedisOption(conf))
	zaplogger := bootstrap.NewLogger(bootstrap.NewLoggerOption(conf))
	postgresql := bootstrap.NewPostgresql(bootstrap.NewPostgresqlOption(conf))

	// gorm.RunMigrations(postgresql)

	userRepo := gormimplement.NewUserRepository(postgresql)
	messageRepo := gormimplement.NewMessageRepository(postgresql)
	channelRepo := gormimplement.NewChannelRepository(postgresql)
	tokenRepo := redisimplement.NewTokenRepository(redispool, conf.GetDuration("jwt.expiration"))

	userUsecase := usecase.NewUserUsecase(userRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo)
	channelUsecase := usecase.NewChannelUsecase(channelRepo)
	tokenUsecase := usecase.NewTokenUsecase(tokenRepo)

	response := response.NewJSONResponse(zaplogger)

	hub := connection.NewHub()
	go hub.Run()

	userController := controller.NewUserController(response, tokenUsecase, userUsecase)
	messageController := controller.NewMessageController(response, messageUsecase)
	channelController := controller.NewChannelController(response, channelUsecase)
	websocketController := controller.NewWebSocketController(response, hub, tokenUsecase, messageUsecase)

	jwtMiddle := jwtMid.NewJWT(response, tokenUsecase)
	corsMiddle := corsMid.NewCORS(corsMid.NewOption((conf)))
	secureHeaderMiddle := secureheader.NewSecureHeader()
	redisRateMiddle := redisrate.NewRateLimiter(redispool, zaplogger, redisrate.NewOption(conf))

	userRouter := router.NewUserRouter(userController, jwtMiddle)
	messageRouter := router.NewMessageRouter(messageController, jwtMiddle)
	channelRouter := router.NewChannelRouter(channelController)
	websocketRouter := router.NewWebSocketRouter(websocketController)

	server := bootstrap.NewServer(
		[]router.IRoute{
			userRouter,
			messageRouter,
			channelRouter,
			websocketRouter,
		},
		[]middleware.IMiddleware{
			corsMiddle,
			secureHeaderMiddle,
		},
	)

	bootstrap.Run(server, bootstrap.NewServerOption(conf))
}
