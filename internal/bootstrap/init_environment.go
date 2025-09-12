package bootstrap

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Environment string

const (
	DEV        Environment = "dev"
	TEST       Environment = "test"
	STAGE      Environment = "stage"
	PRODUCTION Environment = "production"
)

type AppConfig struct {
	env        Environment
	Logger     loggerOption
	Server     serverOption
	Redis      redisoption
	Postgresql postgresqlOption
	Promethus  promethusOption
}

func SetEnvironment(v *viper.Viper, env Environment) (*AppConfig, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = fmt.Sprintf("config/config.%s.yaml", env)
	}
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	appConfig := &AppConfig{
		env:        env,
		Logger:     NewLoggerOption(v),
		Server:     NewServerOption(v),
		Redis:      NewRedisOption(v),
		Postgresql: NewPostgresqlOption(v),
		Promethus:  NewPromethusOption(v),
	}

	return appConfig, nil
}
