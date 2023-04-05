package task

import (
	"context"
	"encoding/json"

	"github.com/getsentry/sentry-go"
	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/watchdog"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type RunCheckTask struct {
	dog                 *watchdog.WatchDog
	monitorRegionDomain *domain.MonitorRegionDomain
}

type RunCheckTaskPayload struct {
	MonitorRegionID int64
}

func NewRunCheckTask(
	dog *watchdog.WatchDog,
	monitorRegionDomain *domain.MonitorRegionDomain,
) *RunCheckTask {
	return &RunCheckTask{
		dog:                 dog,
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
	r.dog.Launch(ctx, monitorRegionWithAssertions)
	return nil
}
