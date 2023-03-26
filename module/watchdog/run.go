package watchdog

import (
	"context"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

func Run(ctx context.Context, monitor *model.Monitor, region *model.Region) {
	tracingID := pkg.GetTracingID(ctx)

	lgr.Default.Print(tracingID, "Running =>", monitor.URL, "from", region.Name)
}
