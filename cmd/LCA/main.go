package main

import (
	"github.com/wang900115/LCA/internal/adapter/controller"
	response "github.com/wang900115/LCA/internal/adapter/controller/response/json"
	"github.com/wang900115/LCA/internal/adapter/middleware"
	corsMid "github.com/wang900115/LCA/internal/adapter/middleware/cors"
	jwtMid "github.com/wang900115/LCA/internal/adapter/middleware/jwt"
	websocketcore "github.com/wang900115/LCA/internal/adapter/websocket/core"
	"github.com/wang900115/LCA/internal/bootstrap"
	"github.com/wang900115/LCA/internal/implement"

	// redisrate "LCA/internal/adapter/gin/middleware/redis_rate"
	secureheader "github.com/wang900115/LCA/internal/adapter/middleware/secure_header"
	"github.com/wang900115/LCA/internal/adapter/router"
	"github.com/wang900115/LCA/internal/application/usecase"
)

func main() {
	conf := bootstrap.NewConfig()

	redispool := bootstrap.NewRedisPool(bootstrap.NewRedisOption(conf))
	zaplogger := bootstrap.NewLogger(bootstrap.NewLoggerOption(conf))
	postgresql := bootstrap.NewPostgresql(bootstrap.NewPostgresqlOption(conf))

	// gorm.RunMigrations(postgresql)

	userRepo := implement.NewUserRepository(postgresql, redispool)
	messageRepo := implement.NewMessageRepository(postgresql, redispool)
	channelRepo := implement.NewChannelRepository(postgresql, redispool)
	// !todo
	tokenRepo := implement.NewTokenRepository(redispool, conf.GetDuration("jwt.login_expiration"), conf.GetDuration("jwt.join_expiration"), []byte("0"), []byte("0"))

	userUsecase := usecase.NewUserUsecase(userRepo, tokenRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo)
	channelUsecase := usecase.NewChannelUsecase(channelRepo)

	response := response.NewJSONResponse(zaplogger)

	hub := websocketcore.NewHub()
	go hub.Run()

	userController := controller.NewUserController(response, userUsecase, channelUsecase, messageUsecase)
	channelController := controller.NewChannelController(response, channelUsecase)
	websocketController := controller.NewWebSocketController(response, hub)

	authjwtMiddle := jwtMid.NewUSERJWT(response, tokenRepo)
	joinjwtMiddle := jwtMid.NewCHANNELJWT(response, tokenRepo)

	corsMiddle := corsMid.NewCORS(corsMid.NewOption((conf)))
	secureHeaderMiddle := secureheader.NewSecureHeader()
	// redisRateMiddle := redisrate.NewRateLimiter(redispool, zaplogger, redisrate.NewOption(conf))

	userRouter := router.NewUserRouter(userController, authjwtMiddle, joinjwtMiddle)
	channelRouter := router.NewChannelRouter(channelController, authjwtMiddle, joinjwtMiddle)
	// messageRouter := router.NewMessageRouter(messageController)
	websocketRouter := router.NewWebSocketRouter(websocketController)

	server := bootstrap.NewServer(
		[]router.IRoute{
			userRouter,
			// messageRouter,
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
