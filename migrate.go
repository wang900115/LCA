package main

import (
	"log"

	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=wang900115 dbname=LCA sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = db.AutoMigrate(&gormmodel.Channel{}, &gormmodel.User{}, &gormmodel.Message{},
		&gormmodel.ChannelMessage{}, &gormmodel.ChannelUser{}, &gormmodel.UserLogin{}, &gormmodel.UserChannel{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
}
