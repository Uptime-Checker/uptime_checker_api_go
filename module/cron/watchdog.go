package cron

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/sourcegraph/conc/iter"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

func (c *Cron) startWatchDog() {
	ctx := pkg.NewTracingID(context.Background())
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Print(tid, 1, "running startWatchDog")

	monitors, err := c.monitorDomain.ListMonitorsToRun(
		ctx,
		-watchDogCheckCronFromAndToInSeconds,
		+watchDogCheckCronFromAndToInSeconds,
	)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	lgr.Print(tid, 2, "number of monitors to run:", len(monitors))

	iterator := iter.Iterator[model.Monitor]{
		MaxGoroutines: watchDogCheckMaxGoroutine,
	}
	iterator.ForEach(monitors, func(monitor *model.Monitor) {
		cachedMonitorNextCheckAt := cache.GetMonitorToRun(monitor.ID)
		if cachedMonitorNextCheckAt == nil {
			// schedule the monitor
		}
	})
}
