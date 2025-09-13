package bootstrap

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"github.com/wang900115/LCA/internal/task"
)

type schedularOption struct {
	Location *time.Location
	Delay    time.Duration
}

func defaultSchedularOption() schedularOption {
	return schedularOption{
		Location: time.Local,
		Delay:    time.Second,
	}
}

func NewSchedularOption(conf *viper.Viper) schedularOption {
	defaultOptions := defaultSchedularOption()
	if conf.IsSet("scheduler.location") {
		if l, err := time.LoadLocation(conf.GetString("scheduler.location")); err == nil {
			defaultOptions.Location = l
		}
	}
	if conf.IsSet("scheduler.delay") {
		defaultOptions.Delay = conf.GetDuration("scheduler.delay")
	}
	return defaultOptions
}

type Scheduler struct {
	jobs []task.IJob
}

func NewScheduler(jobs []task.IJob) *Scheduler {
	return &Scheduler{jobs: jobs}
}

func (s *Scheduler) Run(option schedularOption) *gocron.Scheduler {
	g, err := gocron.NewScheduler(gocron.WithLocation(option.Location))
	if err != nil {
		log.Fatalf("[SYSTEM] Scheduler init error: %s", err.Error())
	}
	for _, job := range s.jobs {
		job.SetUp(g)
	}

	return &g
}
