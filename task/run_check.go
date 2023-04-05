package task

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/watchdog"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
)

type RunCheckTask struct {
	dog                 *watchdog.WatchDog
	monitorDomain       *domain.MonitorDomain
	monitorRegionDomain *domain.MonitorRegionDomain
}

type RunCheckTaskPayload struct {
	MonitorRegionID int64
}

func NewRunCheckTask(
	dog *watchdog.WatchDog,
	monitorDomain *domain.MonitorDomain,
	monitorRegionDomain *domain.MonitorRegionDomain,
) *RunCheckTask {
	return &RunCheckTask{
		dog:                 dog,
		monitorDomain:       monitorDomain,
		monitorRegionDomain: monitorRegionDomain,
	}
}

func (r RunCheckTask) Do(ctx context.Context, job *gue.Job) error {
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Print(tid, 1, "running RunCheckTask")

	var body RunCheckTaskPayload
	if err := json.Unmarshal(job.Args, &body); err != nil {
		return err
	}

	monitorRegionWithAssertions, err := r.monitorRegionDomain.GetWithAllAssoc(ctx, body.MonitorRegionID)
	if err != nil {
		sentry.CaptureException(err)
	}
	go r.dog.Launch(ctx, monitorRegionWithAssertions)

	now := times.Now()
	monitor := monitorRegionWithAssertions.Monitor
	nextCheckAt := now.Add(time.Duration(*monitor.Interval) * time.Second)

	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := r.monitorDomain.UpdateNextCheckAt(ctx, tx, monitor.ID, &now, &nextCheckAt)
		return err
	}); err != nil {
		sentry.CaptureException(err)
	}
	return nil
}
