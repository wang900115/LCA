package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wang900115/LCA/internal/adapter/controller"
	"github.com/wang900115/LCA/internal/adapter/route"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/bootstrap"
	"github.com/wang900115/LCA/pkg/common/middleware"
	middlewareCORS "github.com/wang900115/LCA/pkg/common/middleware/cors"
	middlewareJWT "github.com/wang900115/LCA/pkg/common/middleware/jwt"
	middlewareLOGGER "github.com/wang900115/LCA/pkg/common/middleware/logger"
	middlewarePermission "github.com/wang900115/LCA/pkg/common/middleware/role"
	middlewareSecure "github.com/wang900115/LCA/pkg/common/middleware/secure_header"
	"github.com/wang900115/LCA/pkg/common/router"

	response "github.com/wang900115/LCA/pkg/common/response/json"
	"github.com/wang900115/LCA/pkg/implement"
)

func main() {
	conf := bootstrap.NewConfig()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET_KEY")

	dbGroup := bootstrap.NewDBGroup(conf)
	redisGroup := bootstrap.NewRedisGroup(conf)

	syslogger := bootstrap.NewLogger(bootstrap.NewLoggerOption(conf))

	// !todo 新增一個 check ticker (corn) 來準時 健康檢查 如果有發生問題 -> failover
	cr := implement.NewChannelReadRepository(dbGroup, redisGroup, syslogger)
	cw := implement.NewChannelWriteRepository(dbGroup, redisGroup, syslogger)

	mr := implement.NewMessageReadRepository(dbGroup, redisGroup, syslogger)
	mw := implement.NewMessageWriteRepository(dbGroup, redisGroup, syslogger)

	ur := implement.NewUserReadRepository(dbGroup, redisGroup, syslogger)
	uw := implement.NewUserWriteRepository(dbGroup, redisGroup, syslogger)

	// !todo 要分成 CQRS
	tu := implement.NewTokenAuthRepository(redisGroup.Write, syslogger)

	channel := usecase.NewChannelUsecase(&cr, &cw)
	message := usecase.NewMessageUsecase(&mr, &mw)
	user := usecase.NewUserUsecase(&ur, &uw, &tu, secretKey)

	resp := response.NewJSONResponse(syslogger)

	channelCon := controller.NewChannelController(channel, resp)
	channelMessageCon := controller.NewChannelMessageController(message, resp)
	channelUserCon := controller.NewChannelUserController(channel, resp)
	messageCon := controller.NewMessageController(message, resp)
	userCon := controller.NewUserController(user, resp)
	userChannelCon := controller.NewUserChannelController(channel, resp)
	userChannelMessageCon := controller.NewUserChannelMessageController(message, resp)

	midCORS := middlewareCORS.NewCORS(middlewareCORS.NewOption(conf))
	midJWT := middlewareJWT.NewJWT(resp, &tu, secretKey)
	midRole := middlewarePermission.NewPermission(resp, &tu, secretKey)
	midLog := middlewareLOGGER.NewLogger(syslogger)
	// midRate := middlewareRate.NewRateLimiter(middlewareRate.NewOption(conf))
	midSecure := middlewareSecure.NewSecureHeader()

	userRoute := route.NewUserRouter(userCon, midJWT)
	userChannelRoute := route.NewUserChannelRouter(userChannelCon, midJWT)
	userChannelMessageRoute := route.NewUserChannelMessageRouter(userChannelMessageCon, midJWT)
	channelRoute := route.NewChannelRouter(channelCon)
	channelUserRoute := route.NewChannelUserRouter(channelUserCon, midJWT, midRole)
	channelMessageRoute := route.NewChannelMessageRouter(channelMessageCon)
	messageRoute := route.NewMessageRouter(messageCon, midJWT)

	server := bootstrap.NewServer(
		[]router.IRoute{
			userRoute,
			userChannelRoute,
			userChannelMessageRoute,
			channelRoute,
			channelUserRoute,
			channelMessageRoute,
			messageRoute,
		},
		[]middleware.IMiddleware{
			midCORS,
			midSecure,
			midLog,
		},
	)

	bootstrap.Run(server, bootstrap.NewServerOption(conf))

}
