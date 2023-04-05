package task

import (
	"context"
	"encoding/json"

	"github.com/getsentry/sentry-go"
	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/watchdog"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type StartMonitorTask struct {
	dog              *watchdog.WatchDog
	monitorDomain    *domain.MonitorDomain
	regionDomain     *domain.RegionDomain
	assertionsDomain *domain.AssertionDomain
}

type StartMonitorTaskPayload struct {
	MonitorID int64
}

func NewStartMonitorTask(
	dog *watchdog.WatchDog,
	monitorDomain *domain.MonitorDomain,
	regionDomain *domain.RegionDomain,
	assertionsDomain *domain.AssertionDomain,
) *StartMonitorTask {
	return &StartMonitorTask{
		dog:              dog,
		monitorDomain:    monitorDomain,
		regionDomain:     regionDomain,
		assertionsDomain: assertionsDomain,
	}
}

func (s StartMonitorTask) Do(ctx context.Context, job *gue.Job) error {
	ctx = pkg.NewTracingID(ctx)
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Print(tid, 1, "running StartMonitorTask")

	var body StartMonitorTaskPayload
	if err := json.Unmarshal(job.Args, &body); err != nil {
		return err
	}

	if config.Region == nil {
		region, err := s.regionDomain.Get(ctx, config.App.FlyRegion)
		if err != nil {
			sentry.CaptureException(err)
			return nil
		}
		config.Region = region
	}
	monitor, assertions, err := s.getResources(ctx, body.MonitorID)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	s.dog.Start(ctx, monitor, config.Region, assertions)

	return nil
}

func (s StartMonitorTask) getResources(
	ctx context.Context,
	monitorID int64,
) (*model.Monitor, []model.Assertion, error) {
	monitor, err := s.monitorDomain.Get(ctx, monitorID)
	if err != nil {
		return nil, nil, err
	}

	assertions, err := s.assertionsDomain.ListAssertions(ctx, monitorID)
	return monitor, assertions, err
}
