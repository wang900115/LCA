package infrastructurejob

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostgresqlJob struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewPostgresqlJob(logger *zap.Logger, db *gorm.DB) *PostgresqlJob {
	return &PostgresqlJob{logger: logger, db: db}
}

func (pj *PostgresqlJob) SetUp(s gocron.Scheduler) {
	_, err := s.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(pj.Health))
	if err != nil {
		pj.logger.Error(err.Error(), zap.String("action", "[setup]infrastruction-postgresql-health"))
	}
}

func (pj *PostgresqlJob) Health() {
	sqlDB, err := pj.db.DB()
	if err != nil {
		pj.logger.Error(err.Error(), zap.String("action", "infrastruction-postgresql-health"))
	}
	if err := sqlDB.Ping(); err != nil {
		pj.logger.Error(err.Error(), zap.String("action", "infrastruction-postgresql-health"))
	}
}
