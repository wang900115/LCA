package bootstrap

import (
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type kafkaConsumerOption struct {
	Brokers  []string
	GroupID  string
	Username string
	Password string
}

func NewKafkaConsumerOption(conf *viper.Viper) kafkaConsumerOption {
	return kafkaConsumerOption{
		Brokers:  conf.GetStringSlice("kafka.writer.brokers"),
		GroupID:  conf.GetString("kafka.writer.group_id"),
		Username: conf.GetString("kafka.writer.username"),
		Password: conf.GetString("kafka.writer.password"),
	}
}

func NewKafkaConsumer(option kafkaProducerOption, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  option.Brokers,
		GroupID:  option.GroupID,
		Topic:    topic,
		MaxBytes: 10e6,
		MaxWait:  time.Second,
	})
}
