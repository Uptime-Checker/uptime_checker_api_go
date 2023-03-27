package watchdog

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type WatchDog struct {
	checkDomain *domain.CheckDomain
}

func NewWatchDog(checkDomain *domain.CheckDomain) *WatchDog {
	return &WatchDog{checkDomain: checkDomain}
}

func (c *WatchDog) run(ctx context.Context, tx *sql.Tx, monitor *model.Monitor, region *model.Region) error {
	tracingID := pkg.GetTracingID(ctx)

	lgr.Print(tracingID, 1, "running =>", monitor.URL, "from", region.Name)

	check := &model.Check{
		Success:        false,
		RegionID:       &region.ID,
		MonitorID:      &monitor.ID,
		OrganizationID: monitor.OrganizationID,
	}
	check, err := c.checkDomain.Create(ctx, tx, check)
	if err != nil {
		return err
	}
	lgr.Print(tracingID, 2, "created check", check.ID)

	var headers *map[string]string
	if monitor.Headers != nil {
		if err := json.Unmarshal([]byte(*monitor.Headers), headers); err != nil {
			sentry.CaptureException(err)
		}
	}

	var bodyFormat *resource.MonitorBodyFormat
	if monitor.BodyFormat != nil {
		resourceBodyFormat := resource.MonitorBodyFormat(*monitor.BodyFormat)
		bodyFormat = &resourceBodyFormat
	}

	method := resource.GetMonitorMethod(*monitor.Method)
	hitResponse, hitError := c.hit(
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
	}

	// Update the check
	// Schedule next check
	// Send for alarm if needed
	return nil
}

func (c *WatchDog) gateCheck() {
	// If active subscription,
	// If monitor on,
	// If run too quickly
}
