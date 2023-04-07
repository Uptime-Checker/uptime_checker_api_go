package task

import (
	"context"
	"encoding/json"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
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
	return r.process(ctx, job.Args)
}

func (r RunCheckTask) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return r.process(ctx, t.Payload())
}

func (r RunCheckTask) process(ctx context.Context, payload []byte) error {
	ctx = pkg.NewTracingID(ctx)
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Print(tid, 1, "running RunCheckTask")

	var body RunCheckTaskPayload
	if err := json.Unmarshal(payload, &body); err != nil {
		return err
	}

	if cache.GetMonitorRegionRunning(body.MonitorRegionID) != nil {
		lgr.Print(tid, 2, "monitor region already running", body.MonitorRegionID)
		return nil
	}
	cache.SetMonitorRegionRunning(body.MonitorRegionID)
	monitorRegionWithAssertions, err := r.monitorRegionDomain.GetWithAllAssoc(ctx, body.MonitorRegionID)
	if err != nil {
		sentry.CaptureException(err)
	}
	go r.dog.Launch(ctx, monitorRegionWithAssertions)

	now := times.Now()
	monitor := monitorRegionWithAssertions.Monitor
	nextCheckAt := now.Add(time.Duration(*monitor.Interval) * time.Second)

	lgr.Print(tid, 2, "updating next check for monitor", monitor.ID, "==>", times.Format(nextCheckAt))
	_, err = r.monitorDomain.UpdateNextCheckAt(ctx, monitor.ID, &now, &nextCheckAt)
	if err != nil {
		sentry.CaptureException(err)
	}
	return nil
}
