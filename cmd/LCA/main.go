package main

import (
	"github.com/spf13/viper"
	"github.com/wang900115/LCA/internal/adapter/controller"
	response "github.com/wang900115/LCA/internal/adapter/controller/response/json"
	"github.com/wang900115/LCA/internal/adapter/middleware"
	rbacMid "github.com/wang900115/LCA/internal/adapter/middleware/casbin"
	corsMid "github.com/wang900115/LCA/internal/adapter/middleware/cors"
	jwtMid "github.com/wang900115/LCA/internal/adapter/middleware/jwt"
	websocketcore "github.com/wang900115/LCA/internal/adapter/websocket/core"
	"github.com/wang900115/LCA/internal/bootstrap"
	"github.com/wang900115/LCA/internal/implement"
	"github.com/wang900115/LCA/internal/task"
	infrastructurejob "github.com/wang900115/LCA/internal/task/infrastructure-job"

	// redisrate "LCA/internal/adapter/gin/middleware/redis_rate"
	secureheader "github.com/wang900115/LCA/internal/adapter/middleware/secure_header"
	"github.com/wang900115/LCA/internal/adapter/router"
	"github.com/wang900115/LCA/internal/application/usecase"
)

func main() {
	conf := viper.New()
	appOptions, err := bootstrap.SetEnvironment(conf, "dev")
	if err != nil {
		panic(err)
	}
	redispool := bootstrap.NewRedisPool(appOptions.Redis)
	zaplogger := bootstrap.NewLogger(appOptions.Logger)
	postgresql := bootstrap.NewPostgresql(appOptions.Postgresql)
	casbin := bootstrap.NewCasbin(postgresql, appOptions.Casbin)

	job1 := infrastructurejob.NewPostgresqlJob(zaplogger, postgresql)
	job2 := infrastructurejob.NewRedisJob(zaplogger, redispool)

	scheduler := bootstrap.NewScheduler(
		[]task.IJob{
			job1,
			job2,
		})
	// promethus := bootstrap.NewPromethus(appOptions.Promethus)

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
	websocketController := controller.NewWebSocketController(response, userUsecase, hub)

	rabcMiddle := rbacMid.NewCASBIN(response, casbin)
	authjwtMiddle := jwtMid.NewUSERJWT(response, tokenRepo)
	joinjwtMiddle := jwtMid.NewCHANNELJWT(response, tokenRepo)

	corsMiddle := corsMid.NewCORS(corsMid.NewOption((conf)))
	secureHeaderMiddle := secureheader.NewSecureHeader()
	// redisRateMiddle := redisrate.NewRateLimiter(redispool, zaplogger, redisrate.NewOption(conf))

	userRouter := router.NewUserRouter(userController, authjwtMiddle, joinjwtMiddle, rabcMiddle)
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

	srv := server.Run(appOptions.Server)
	sch := scheduler.Run(appOptions.Gocron)

	bootstrap.Run(appOptions.Server.CancelTimeout, srv, *sch)
}
