package watchdog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/samber/lo"

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
	checkDomain              *domain.CheckDomain
	regionDomain             *domain.RegionDomain
	monitorDomain            *domain.MonitorDomain
	monitorRegionDomain      *domain.MonitorRegionDomain
	monitorStatusDomain      *domain.MonitorStatusDomain
	monitorIntegrationDomain *domain.MonitorIntegrationDomain
	alarmDomain              *domain.AlarmDomain
	alarmChannelDomain       *domain.AlarmChannelDomain

	checkService         *service.CheckService
	monitorService       *service.MonitorService
	monitorRegionService *service.MonitorRegionService
	errorLogService      *service.ErrorLogService
	dailyReportService   *service.DailyReportService
	alarmPolicyService   *service.AlarmPolicyService
}

func NewWatchDog(
	checkDomain *domain.CheckDomain,
	regionDomain *domain.RegionDomain,
	monitorDomain *domain.MonitorDomain,
	monitorRegionDomain *domain.MonitorRegionDomain,
	monitorStatusDomain *domain.MonitorStatusDomain,
	monitorIntegrationDomain *domain.MonitorIntegrationDomain,
	alarmDomain *domain.AlarmDomain,
	alarmChannelDomain *domain.AlarmChannelDomain,
	checkService *service.CheckService,
	monitorService *service.MonitorService,
	monitorRegionService *service.MonitorRegionService,
	errorLogService *service.ErrorLogService,
	dailyReportService *service.DailyReportService,
	alarmPolicyService *service.AlarmPolicyService,
) *WatchDog {
	return &WatchDog{
		checkDomain:              checkDomain,
		regionDomain:             regionDomain,
		monitorDomain:            monitorDomain,
		monitorRegionDomain:      monitorRegionDomain,
		monitorStatusDomain:      monitorStatusDomain,
		monitorIntegrationDomain: monitorIntegrationDomain,
		alarmDomain:              alarmDomain,
		alarmChannelDomain:       alarmChannelDomain,
		checkService:             checkService,
		monitorService:           monitorService,
		monitorRegionService:     monitorRegionService,
		errorLogService:          errorLogService,
		dailyReportService:       dailyReportService,
		alarmPolicyService:       alarmPolicyService,
	}
}

// Launch is run by the cron
func (w *WatchDog) Launch(
	ctx context.Context,
	monitorRegionWithAssertions *pkg.MonitorRegionWithAssertions,
) {
	tracingID := pkg.GetTracingID(ctx)
	monitor := monitorRegionWithAssertions.Monitor

	check, hitResponse, hitErr, err := w.fly(ctx, monitor.Monitor, monitorRegionWithAssertions.Region)
	if err != nil {
		sentry.CaptureException(err)
	}

	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
		check, errorLog, err := w.run(
			ctx, tx, check, monitor.Monitor, monitor.Assertions, hitResponse, hitErr,
		)
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 1, "check ran, successful:", check.Success,
			"duration:", fmt.Sprintf("%dms", *check.Duration))
		// Insert to the daily report
		dailyReport, err := w.dailyReportService.Add(ctx, tx, monitor.ID, monitor.OrganizationID, check.Success)
		if err != nil {
			return errors.Newf("daily report add failed, err: %w", err)
		}
		// Send for verification and alarm
		return w.verify(
			ctx, tx, check, errorLog, monitor.Monitor, monitorRegionWithAssertions.MonitorRegion, dailyReport,
		)
	}); err != nil {
		sentry.CaptureException(err)
	}
}

// Start is run by the controller
func (w *WatchDog) Start(
	ctx context.Context,
	monitor *model.Monitor,
	region *model.Region,
	assertions []model.Assertion,
) {
	if err := w.startMonitor(ctx, monitor, region, assertions); err != nil {
		sentry.CaptureException(err)
	}
}

func (w *WatchDog) startMonitor(
	ctx context.Context,
	monitor *model.Monitor,
	region *model.Region,
	assertions []model.Assertion,
) error {
	tracingID := pkg.GetTracingID(ctx)
	if region == nil {
		region, err := w.regionDomain.Get(ctx, config.App.FlyRegion)
		if err != nil {
			return errors.Newf("failed to get region %s, err: %w", config.App.FlyRegion, err)
		}
		config.Region = region
	}
	check, hitResponse, hitErr, err := w.fly(ctx, monitor, config.Region)
	if err != nil {
		return err
	}

	return infra.Transaction(ctx, func(tx *sql.Tx) error {
		check, _, err := w.run(
			ctx, tx, check, monitor, assertions, hitResponse, hitErr,
		)
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 1, "check ran, successful:", check.Success,
			"duration:", fmt.Sprintf("%dms", *check.Duration))
		if check.Success {
			monitor, err := w.monitorService.StartOn(ctx, tx, monitor)
			if err != nil {
				return err
			}
			lgr.Print(tracingID, 2, "starting monitor for", monitor.URL)
			_, err = w.monitorRegionService.FirstOrCreate(ctx, tx, monitor.ID, config.Region.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// fly hits the target
func (w *WatchDog) fly(
	ctx context.Context,
	monitor *model.Monitor,
	region *model.Region,
) (*model.Check, *HitResponse, *HitErr, error) {
	tracingID := pkg.GetTracingID(ctx)
	lgr.Print(tracingID, 1, "running =>", monitor.URL, "from", region.Name)

	check := &model.Check{
		Success:        false,
		RegionID:       region.ID,
		MonitorID:      monitor.ID,
		OrganizationID: monitor.OrganizationID,
	}
	check, err := w.checkDomain.Create(ctx, check)
	if err != nil {
		return nil, nil, nil, errors.Newf("failed to create check, monitor %d, err: %w", monitor.ID, err)
	}
	lgr.Print(tracingID, 2, "created check", check.ID)

	var headers map[string]string
	if monitor.Headers != nil {
		if err := json.Unmarshal([]byte(*monitor.Headers), &headers); err != nil {
			return nil, nil, nil, errors.Newf("could not unmarshal monior headers %d, err: %w", monitor.ID, err)
		}
	}

	method := resource.GetMonitorMethod(*monitor.Method)
	hitResponse, hitError := w.Hit(
		ctx,
		monitor.URL,
		method,
		monitor.Body,
		monitor.Username,
		monitor.Password,
		resource.MonitorBodyFormat(monitor.BodyFormat),
		&headers,
		monitor.Timeout,
		monitor.FollowRedirects,
	)
	return check, hitResponse, hitError, err
}

// run comes after flying
func (w *WatchDog) run(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	monitor *model.Monitor,
	assertions []model.Assertion,
	hitResponse *HitResponse, hitError *HitErr,
) (*model.Check, *model.ErrorLog, error) {
	tracingID := pkg.GetTracingID(ctx)

	var err error
	var errorLog *model.ErrorLog
	method := resource.GetMonitorMethod(*monitor.Method)

	if hitResponse == nil && hitError != nil {
		lgr.Print(tracingID, 1, "hit request failed", method, monitor.URL)
		// Create error log
		errorLog, err = w.errorLogService.Create(ctx, tx, monitor.ID, check.ID, nil, &hitError.Text, hitError.Type)
		if err != nil {
			return check, nil, errors.Newf("failed to create error log, monitor %d, err: %w", monitor.ID, err)
		}
	} else {
		checkSuccess := true
		// Assertion test
		var failedAssertion *model.Assertion
		if hitError == nil {
			for i, assertion := range assertions {
				if pass := w.Assert(
					assertion.Source, assertion.Property, assertion.Comparison, *assertion.Value, *hitResponse,
				); !pass {
					failedAssertion = &assertions[i]
					break
				}
			}
		} else {
			checkSuccess = false
			// Create error log
			errorLog, err = w.errorLogService.Create(ctx, tx, monitor.ID, check.ID, nil,
				&hitError.Text, hitError.Type)
			if err != nil {
				return check, nil, errors.Newf("failed to create error log, check %d, err: %w", check.ID, err)
			}
		}

		if failedAssertion != nil {
			checkSuccess = false
			// Create error log
			errorLog, err = w.errorLogService.Create(ctx, tx, monitor.ID, check.ID, lo.ToPtr(failedAssertion.ID),
				nil, resource.ErrorLogTypeAssertionFailure)
			if err != nil {
				return check, nil, errors.Newf("failed to create error log, check %d, err: %w", check.ID, err)
			}
		}

		// update the check
		check, err = w.checkService.Update(ctx, tx, check, checkSuccess, hitResponse.Duration, hitResponse.Size,
			hitResponse.StatusCode, hitResponse.ContentType, hitResponse.Body, hitResponse.Headers, hitResponse.Traces)
		if err != nil {
			return check, nil, errors.Newf("failed to update check %d, err: %w", check.ID, err)
		}
	}
	return check, errorLog, nil
}

func (w *WatchDog) gateCheck() {
	// If active subscription,
	// If monitor on,
	// If run too quickly
}
