package bootstrap

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type kafkaOption struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewKafkaOption(conf *viper.Viper) kafkaOption {
	return kafkaOption{
		Host:     conf.GetString("kafka.Host"),
		Port:     conf.GetString("kafka.Port"),
		Username: conf.GetString("kafka.Username"),
		Password: conf.GetString("kafka.Password"),
	}
}

func NewKafka(option kafkaOption) *kafka.Conn {
	conn, err := kafka.Dial("tcp", option.Host+option.Port)
	if err != nil {
		panic(err)
	}
	return conn
}
