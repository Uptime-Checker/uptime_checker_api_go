package watchdog

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
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

func (c *WatchDog) Run(ctx context.Context, tx *sql.Tx, monitor *model.Monitor, region *model.Region) error {
	tracingID := pkg.GetTracingID(ctx)

	lgr.Default.Print(tracingID, 1, "running =>", monitor.URL, "from", region.Name)

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
	lgr.Default.Print(tracingID, 2, "created check", check.ID)
	return nil
}
