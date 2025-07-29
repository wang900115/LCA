package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
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
	response "github.com/wang900115/LCA/pkg/common/response/json"
	"github.com/wang900115/LCA/pkg/common/router"
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
	applogger := bootstrap.NewLogger(bootstrap.NewLoggerOption(conf))

	ctx := context.Background()

	c := cron.New()

	if err := c.AddFunc("/30 * * * * *", func() {
		redisGroup.HeadlthCheck(ctx)
	}); err != nil {
		log.Fatalf("Failed to add cron redis-health-check func: %v", err)
	}

	if err := c.AddFunc("/60 * * * * *", func() {
		dbGroup.HeadlthCheck(ctx)
	}); err != nil {
		log.Fatalf("Failed to add cron postgre-health-check func: %v", err)
	}

	c.Start()

	cr := implement.NewChannelReadRepository(dbGroup, redisGroup, syslogger)
	cw := implement.NewChannelWriteRepository(dbGroup, redisGroup, syslogger)

	mr := implement.NewMessageReadRepository(dbGroup, redisGroup, syslogger)
	mw := implement.NewMessageWriteRepository(dbGroup, redisGroup, syslogger)

	ur := implement.NewUserReadRepository(dbGroup, redisGroup, syslogger)
	uw := implement.NewUserWriteRepository(dbGroup, redisGroup, syslogger)

	tu := implement.NewTokenAuthRepository(redisGroup, syslogger)

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
	midLog := middlewareLOGGER.NewLogger(applogger)
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

	bootstrap.Run(server, bootstrap.NewServerOption(conf), c)

}
