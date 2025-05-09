package gorm

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Option struct {
	User     string
	Password string
	Host     string
	Port     string
	DBname   string
	SSLMode  string
}

func NewOption(conf *viper.Viper) Option {
	return Option{
		User:     conf.GetString("postgresql.user"),
		Password: conf.GetString("postgresql.password"),
		Host:     conf.GetString("postgresql.host"),
		Port:     conf.GetString("postgresql.port"),
		DBname:   conf.GetString("postgresql.dbname"),
		SSLMode:  conf.GetString("postgresql.sslmode"),
	}
}

func NewPostgresql(option Option) *gorm.DB {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		option.User, option.Password, option.Host, option.Port, option.DBname, option.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	return db
}
