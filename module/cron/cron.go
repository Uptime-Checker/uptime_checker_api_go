package cron

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

var s *gocron.Scheduler

type Task interface {
	Do()
}

// JobName type
type JobName string

const (
	JobNameSyncStripeProducts JobName = "SYNC_STRIPE_PRODUCTS"
)

type Cron struct {
	jobDomain *domain.JobDomain

	syncProductsTask *task.SyncProductsTask
}

func NewCron(jobDomain *domain.JobDomain, syncProductsTask *task.SyncProductsTask) *Cron {
	return &Cron{jobDomain: jobDomain, syncProductsTask: syncProductsTask}
}

func (c *Cron) Start() error {
	s = gocron.NewScheduler(time.UTC)
	now := times.Now()

	random := pkg.RandomNumber(60, 120)
	_, err := s.Every(30).Second().StartAt(now.Add(time.Second * time.Duration(random))).Do(c.checkAndRun)
	if err != nil {
		return err
	}

	s.StartAsync()
	return nil
}

func (c *Cron) checkAndRun() {
	ctx := context.Background()
	lgr.Default.Print("Running cron check")

	jobsToRun, err := c.jobDomain.ListJobsToRun(ctx, -20, 20)
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	for i, job := range jobsToRun {
		if job.Name == string(JobNameSyncStripeProducts) {
			go runTask(ctx, c.jobDomain, *c.syncProductsTask, jobsToRun[i])
		}
	}
}

func runTask[T Task](ctx context.Context, jobDomain *domain.JobDomain, task T, job model.Job) {
	now := times.Now()
	nextRunAt := now
	if *job.Recurring {
		nextRunAt = now.Add(time.Minute * time.Duration(*job.Interval))
	}

	_, err := jobDomain.UpdateRunning(ctx, job.ID, &now, &nextRunAt, resource.JobStatusRunning)
	if err != nil {
		sentry.CaptureException(err)
	}
	task.Do()

	if *job.Recurring {
		_, err = jobDomain.UpdateStatus(ctx, job.ID, resource.JobStatusScheduled)
		if err != nil {
			sentry.CaptureException(err)
		}
	}
}
