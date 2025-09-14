package infrastructurejob

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaJob struct {
	logger *zap.Logger
	kafka  *kafka.Conn
}

func NewKafkaJob(logger *zap.Logger, kafka *kafka.Conn) *KafkaJob {
	return &KafkaJob{logger: logger, kafka: kafka}
}

func (kj *KafkaJob) SetUp(s gocron.Scheduler) {
	_, err := s.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(kj.Health))
	if err != nil {
		kj.logger.Error(err.Error(), zap.String("action", "[setup]infrastruction-kafka-health"))
	}
}

func (kj *KafkaJob) Health() {
	_, err := kj.kafka.Brokers()
	if err != nil {
		kj.logger.Error(err.Error(), zap.String("action", "infrastruction-kafka-health-broker"))
	}
	_, err = kj.kafka.Controller()
	if err != nil {
		kj.logger.Error(err.Error(), zap.String("action", "infrastruction-kafka-health-controller"))
	}
}
