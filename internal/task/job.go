package task

import "github.com/go-co-op/gocron/v2"

type IJob interface {
	SetUp(s gocron.Scheduler)
}
