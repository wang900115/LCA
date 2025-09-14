package bootstrap

import (
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type kafkaProducerOption struct {
	Brokers  []string
	GroupID  string
	Username string
	Password string
}

func NewKafkaProducerOption(conf *viper.Viper) kafkaProducerOption {
	return kafkaProducerOption{
		Brokers:  conf.GetStringSlice("kafka.writer.brokers"),
		GroupID:  conf.GetString("kafka.writer.group_id"),
		Username: conf.GetString("kafka.writer.username"),
		Password: conf.GetString("kafka.writer.password"),
	}
}

func NewKafkaProducer(option kafkaProducerOption, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  option.Brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}
