package cron

import (
	"context"
	"database/sql"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

const (
	checkCronFromAndToInSeconds         = 20
	watchDogCheckCronFromAndToInSeconds = 4
	watchDogCheckMaxGoroutine           = 100
)

var s *gocron.Scheduler

type Task interface {
	Do(ctx context.Context)
}

// JobName type
type JobName string

const (
	JobNameSyncStripeProducts JobName = "SYNC_STRIPE_PRODUCTS"
	JobNameCheckWatchdog      JobName = "CHECK_WATCHDOG"
)

type Cron struct {
	jobDomain           *domain.JobDomain
	regionDomain        *domain.RegionDomain
	monitorDomain       *domain.MonitorDomain
	monitorRegionDomain *domain.MonitorRegionDomain

	propertyService *service.PropertyService

	syncProductsTask *task.SyncProductsTask
}

func NewCron(
	jobDomain *domain.JobDomain,
	regionDomain *domain.RegionDomain,
	monitorDomain *domain.MonitorDomain,
	monitorRegionDomain *domain.MonitorRegionDomain,
	propertyService *service.PropertyService,
	syncProductsTask *task.SyncProductsTask,
) *Cron {
	return &Cron{
		jobDomain:           jobDomain,
		regionDomain:        regionDomain,
		monitorDomain:       monitorDomain,
		monitorRegionDomain: monitorRegionDomain,
		propertyService:     propertyService,
		syncProductsTask:    syncProductsTask,
	}
}

func Shutdown() {
	s.Stop()
}

func (c *Cron) Start(ctx context.Context) error {
	tracingID := pkg.GetTracingID(ctx)
	now := times.Now()
	s = gocron.NewScheduler(time.UTC)

	// start croner
	_, err := s.Every(constant.CronCheckIntervalInSeconds).Second().
		StartAt(now.Add(time.Second * 5)).
		Do(c.checkAndRun)
	if err != nil {
		return err
	}

	recurringJobs, err := c.jobDomain.ListRecurringJobs(ctx)
	if err != nil {
		sentry.CaptureException(err)
	} else {
		if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
			for _, job := range recurringJobs {
				if job.NextRunAt == nil || times.CompareDate(now, *job.NextRunAt) == constant.Date1AfterDate2 {
					nextRunAt := now.Add(time.Second * time.Duration(*job.Interval))
					lgr.Print("updating the next run at for", job.Name, "to", times.Format(nextRunAt))
					_, err := c.jobDomain.UpdateNextRunAt(ctx, tx, job.ID, &nextRunAt, resource.JobStatusScheduled)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			sentry.CaptureException(err)
		}
	}

	// start the watchdog
	_, err = s.Every(constant.WatchDogCheckIntervalInSeconds).Second().StartAt(now.Add(time.Second * 2)).Do(c.watchDog)
	if err != nil {
		return err
	}

	s.StartAsync()
	lgr.Print(tracingID, "cron started")
	return nil
}

func (c *Cron) checkAndRun() {
	ctx := pkg.NewTracingID(context.Background())
	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, 1, "running cron check")

	// Cron check runs every 30s. We look for jobs that need to be run from last 20s to next 20s from current time
	jobsToRun, err := c.jobDomain.ListJobsToRun(ctx, -checkCronFromAndToInSeconds, checkCronFromAndToInSeconds)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	lgr.Print(tid, 2, "found", len(jobsToRun), "cron jobs to run")

	for i, job := range jobsToRun {
		if job.Name == string(JobNameSyncStripeProducts) {
			go runTask(c.jobDomain, *c.syncProductsTask, jobsToRun[i])
		} else if job.Name == string(JobNameCheckWatchdog) {
			c.stopTheDog(ctx)
		}
	}
}

func runTask[T Task](jobDomain *domain.JobDomain, tsk T, job model.Job) {
	ctx := pkg.NewTracingID(context.Background())
	defer sentry.RecoverWithContext(ctx)
	now := times.Now()
	nextRunAt := now
	if job.Recurring {
		nextRunAt = now.Add(time.Second * time.Duration(*job.Interval))
	}

	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
		_, err := jobDomain.UpdateRunning(ctx, tx, job.ID, &now, &nextRunAt, resource.JobStatusRunning)
		if err != nil {
			return err
		}
		tsk.Do(ctx)
		if job.Recurring {
			_, err = jobDomain.UpdateStatus(ctx, tx, job.ID, resource.JobStatusScheduled)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		sentry.CaptureException(err)
	}
}
