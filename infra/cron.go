package infra

import (
	"time"

	"github.com/go-co-op/gocron"

	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

var s *gocron.Scheduler

func SetupCron(syncProductsTask *task.SyncProductsTask) error {
	s = gocron.NewScheduler(time.UTC)
	now := times.Now()

	_, err := s.Every(63).Minutes().StartAt(now.Add(time.Second * 60)).Do(syncProductsTask.Run)
	if err != nil {
		return err
	}

	s.StartAsync()
	return nil
}
