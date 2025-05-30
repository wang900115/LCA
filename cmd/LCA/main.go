package main

import (
	"github.com/wang900115/LCA/internal/adapter/gin"
	"github.com/wang900115/LCA/internal/adapter/gin/controller"
	response "github.com/wang900115/LCA/internal/adapter/gin/controller/response/json"
	"github.com/wang900115/LCA/internal/adapter/gin/middleware"
	corsMid "github.com/wang900115/LCA/internal/adapter/gin/middleware/cors"
	jwtMid "github.com/wang900115/LCA/internal/adapter/gin/middleware/jwt"
	"github.com/wang900115/LCA/internal/adapter/gorm"

	// redisrate "LCA/internal/adapter/gin/middleware/redis_rate"
	secureheader "github.com/wang900115/LCA/internal/adapter/gin/middleware/secure_header"
	"github.com/wang900115/LCA/internal/adapter/gin/router"
	"github.com/wang900115/LCA/internal/adapter/redispool"
	"github.com/wang900115/LCA/internal/adapter/repository"
	"github.com/wang900115/LCA/internal/adapter/websocket/connection"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/wang900115/LCA/pkg/config"
	"github.com/wang900115/LCA/pkg/logger"
)

func main() {
	conf := config.NewConfig()

	redispool := redispool.NewRedisPool(redispool.NewOption(conf))
	zaplogger := logger.NewZapLogger(logger.NewOption(conf))
	postgresql := gorm.NewPostgresql(gorm.NewOption(conf))

	gorm.RunMigrations(postgresql)

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
	// redisRateMiddle := redisrate.NewRateLimiter(redispool, zaplogger, redisrate.NewOption(conf))

	userRouter := router.NewUserRouter(userController, jwtMiddle)
	messageRouter := router.NewMessageRouter(messageController, jwtMiddle)
	channelRouter := router.NewChannelRouter(channelController)
	websocketRouter := router.NewWebSocketRouter(websocketController)

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
		},
	)

	gin.Run(app, gin.NewOption(conf))
}
