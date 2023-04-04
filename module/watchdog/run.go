package watchdog

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
)

type WatchDog struct {
	checkDomain         *domain.CheckDomain
	regionDomain        *domain.RegionDomain
	assertionDomain     *domain.AssertionDomain
	monitorRegionDomain *domain.MonitorRegionDomain
	monitorStatusDomain *domain.MonitorStatusDomain

	checkService         *service.CheckService
	monitorService       *service.MonitorService
	monitorRegionService *service.MonitorRegionService
}

func NewWatchDog(
	checkDomain *domain.CheckDomain,
	regionDomain *domain.RegionDomain,
	assertionDomain *domain.AssertionDomain,
	monitorRegionDomain *domain.MonitorRegionDomain,
	monitorStatusDomain *domain.MonitorStatusDomain,
	checkService *service.CheckService,
	monitorService *service.MonitorService,
	monitorRegionService *service.MonitorRegionService,
) *WatchDog {
	return &WatchDog{
		checkDomain:          checkDomain,
		regionDomain:         regionDomain,
		assertionDomain:      assertionDomain,
		monitorRegionDomain:  monitorRegionDomain,
		monitorStatusDomain:  monitorStatusDomain,
		checkService:         checkService,
		monitorService:       monitorService,
		monitorRegionService: monitorRegionService,
	}
}

// Launch is run by the cron
func (w *WatchDog) Launch(
	ctx context.Context,
	monitor *model.Monitor,
	region *model.Region,
	monitorRegion *model.MonitorRegion,
) {
}

// Start is run by the controller
func (w *WatchDog) Start(
	ctx context.Context,
	monitor *model.Monitor,
	region *model.Region,
) {
	if err := w.startMonitor(ctx, monitor, region); err != nil {
		sentry.CaptureException(err)
	}
}

func (w *WatchDog) startMonitor(
	ctx context.Context,
	monitor *model.Monitor,
	region *model.Region,
) error {
	tracingID := pkg.GetTracingID(ctx)
	if region == nil {
		region, err := w.regionDomain.Get(ctx, config.App.FlyRegion)
		if err != nil {
			return err
		}
		config.Region = region
	}
	return infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		check, err := w.run(ctx, tx, monitor, config.Region)
		if err != nil {
			return err
		}
		if check.Success {
			monitor, err := w.monitorService.Start(ctx, tx, monitor, true)
			if err != nil {
				return err
			}
			lgr.Print(tracingID, 1, "starting monitor for", monitor.URL)
			_, err = w.monitorRegionService.FirstOrCreate(ctx, tx, monitor.ID, region.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (w *WatchDog) run(
	ctx context.Context,
	tx *sql.Tx,
	monitor *model.Monitor,
	region *model.Region,
) (*model.Check, error) {
	tracingID := pkg.GetTracingID(ctx)

	lgr.Print(tracingID, 1, "running =>", monitor.URL, "from", region.Name)

	check := &model.Check{
		Success:        false,
		RegionID:       &region.ID,
		MonitorID:      &monitor.ID,
		OrganizationID: monitor.OrganizationID,
	}
	check, err := w.checkDomain.Create(ctx, tx, check)
	if err != nil {
		return nil, err
	}
	lgr.Print(tracingID, 2, "created check", check.ID)

	var headers *map[string]string
	if monitor.Headers != nil {
		if err := json.Unmarshal([]byte(*monitor.Headers), headers); err != nil {
			return nil, err
		}
	}

	var bodyFormat *resource.MonitorBodyFormat
	if monitor.BodyFormat != nil {
		resourceBodyFormat := resource.MonitorBodyFormat(*monitor.BodyFormat)
		bodyFormat = &resourceBodyFormat
	}

	method := resource.GetMonitorMethod(*monitor.Method)
	hitResponse, hitError := w.Hit(
		ctx,
		monitor.URL,
		method,
		monitor.Body,
		monitor.Username,
		monitor.Password,
		bodyFormat,
		headers,
		*monitor.Timeout,
		*monitor.FollowRedirects,
	)

	if hitResponse == nil && hitError != nil {
		lgr.Print(tracingID, 1, "hit request failed", method, monitor.URL)
		// Create error log
	} else {
		checkSuccess := true
		// assertion test
		var failedAssertion *model.Assertion
		if hitError == nil {
			assertions, err := w.assertionDomain.ListAssertions(ctx, monitor.ID)
			if err != nil {
				return nil, err
			}
			for i, assertion := range assertions {
				if pass := w.Assert(
					*assertion.Source, assertion.Property, *assertion.Comparison, *assertion.Value, *hitResponse,
				); !pass {
					failedAssertion = &assertions[i]
					break
				}
			}
		} else {
			checkSuccess = false
			// Create error log
		}

		if failedAssertion != nil {
			checkSuccess = false
			// Create error log
		}

		// Update the check
		check, err = w.checkService.Update(ctx, tx, check, checkSuccess, hitResponse.Duration, hitResponse.Size,
			hitResponse.ContentType, hitResponse.Body, hitResponse.Headers, hitResponse.Traces)
		if err != nil {
			return nil, err
		}
	}

	// Schedule next check
	// Send for alarm if needed
	return check, nil
}

func (w *WatchDog) gateCheck() {
	// If active subscription,
	// If monitor on,
	// If run too quickly
}
