package main

import (
	"LCA/internal/adapter/gin"
	"LCA/internal/adapter/gin/controller"
	response "LCA/internal/adapter/gin/controller/response/json"
	"LCA/internal/adapter/gin/middleware"
	corsMid "LCA/internal/adapter/gin/middleware/cors"
	jwtMid "LCA/internal/adapter/gin/middleware/jwt"
	redisrate "LCA/internal/adapter/gin/middleware/redis_rate"
	secureheader "LCA/internal/adapter/gin/middleware/secure_header"
	"LCA/internal/adapter/gin/router"
	"LCA/internal/adapter/redispool"
	"LCA/internal/adapter/repository"
	"LCA/internal/adapter/websocket/connection"
	"LCA/internal/application/usecase"
	"LCA/pkg/config"
	"LCA/pkg/gorm"
	"LCA/pkg/logger"
)

func main() {
	conf := config.NewConfig()

	redispool := redispool.NewRedisPool(redispool.NewOption(conf))
	zaplogger := logger.NewZapLogger(logger.NewOption(conf))
	postgresql := gorm.NewPostgresql(gorm.NewOption(conf))

	userRepo := repository.NewUserRepository(postgresql)
	messageRepo := repository.NewMessageRepository(postgresql)
	channelRepo := repository.NewChannelRepository(postgresql)
	tokenRepo := repository.NewTokenRepository(redispool, conf.GetDuration("jwt.expiration"))

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
	websocketRouter := router.NewWebSocketRouter(websocketController, jwtMiddle)

	app := gin.NewApp(
		[]router.IRoute{
			userRouter,
			messageRouter,
			channelRouter,
			websocketRouter,
		},
		[]middleware.IMiddleware{
			corsMiddle,
			secureHeaderMiddle,
			redisRateMiddle,
		},
	)

	gin.Run(app, gin.NewOption(conf))
}
