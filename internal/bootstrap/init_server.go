package bootstrap

import (
	"net/http"
	"time"

	"github.com/wang900115/LCA/internal/adapter/middleware"
	"github.com/wang900115/LCA/internal/adapter/router"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type serverOption struct {
	RunMode           string
	HTTPPort          string
	CancelTimeout     time.Duration
	ReadHeaderTimeout time.Duration
}

func defaultServerOption() serverOption {
	return serverOption{
		RunMode:           gin.DebugMode,
		HTTPPort:          "8080",
		CancelTimeout:     5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func NewServerOption(conf *viper.Viper) serverOption {
	defaultOptions := defaultServerOption()
	if conf.IsSet("app.run_mode") {
		defaultOptions.RunMode = conf.GetString("app.run_mode")
	}
	if conf.IsSet("server.http_port") {
		defaultOptions.HTTPPort = conf.GetString("server.http_port")
	}
	if conf.IsSet("app.cancel_timeout") {
		defaultOptions.CancelTimeout = conf.GetDuration("app.cancel_timeout")
	}
	if conf.IsSet("app.read_header_timeout") {
		defaultOptions.ReadHeaderTimeout = conf.GetDuration("app.read_header_timeout")
	}
	return defaultOptions
}

type App struct {
	routes      []router.IRoute
	middlewares []middleware.IMiddleware
}

func NewServer(routes []router.IRoute, middlewares []middleware.IMiddleware) *App {
	return &App{
		routes:      routes,
		middlewares: middlewares,
	}
}

func (a *App) Run(option serverOption) *http.Server {
	gin.SetMode(option.RunMode)

	routerEngine := gin.Default()

	for _, middleware := range a.middlewares {
		routerEngine.Use(middleware.Middleware)
	}

	for _, route := range a.routes {
		route.Setup(routerEngine.Group("/api"))
	}

	if option.RunMode == "debug" {
		pprof.Register(routerEngine)
	}

	srv := &http.Server{
		Addr:              ":" + option.HTTPPort,
		Handler:           routerEngine,
		ReadHeaderTimeout: option.ReadHeaderTimeout,
	}

	return srv
}
