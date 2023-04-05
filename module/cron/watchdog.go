package cron

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/sourcegraph/conc/iter"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/task/client"
)

func (c *Cron) watchDog() {
	ctx := pkg.NewTracingID(context.Background())
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Print(tid, 1, "running watchDog")
	if err := c.watchTheDog(ctx, tid); err != nil {
		sentry.CaptureException(err)
	}
}

func (c *Cron) watchTheDog(ctx context.Context, tid string) error {
	if config.Region == nil {
		region, err := c.regionDomain.Get(ctx, config.App.FlyRegion)
		if err != nil {
			return err
		}
		config.Region = region
	}

	monitors, err := c.monitorDomain.ListMonitorsToRun(
		ctx,
		-watchDogCheckCronFromAndToInSeconds,
		+watchDogCheckCronFromAndToInSeconds,
	)
	if err != nil {
		return err
	}
	lgr.Print(tid, 2, "number of monitors to run:", len(monitors))

	iterator := iter.Iterator[model.Monitor]{
		MaxGoroutines: watchDogCheckMaxGoroutine,
	}
	iterator.ForEach(monitors, func(monitor *model.Monitor) {
		cachedMonitorNextCheckAt := cache.GetMonitorToRun(monitor.ID)
		if cachedMonitorNextCheckAt == nil {
			// schedule the monitor
			monitorRegion, err := c.monitorRegionDomain.GetOldestChecked(ctx, monitor.ID, config.Region.ID)
			if err != nil {
				lgr.Error(tid, 3, "failed to get monitor region", err)
				return
			}
			if err := client.RunCheckAsync(ctx, monitorRegion.ID, *monitor.NextCheckAt); err != nil {
				lgr.Error(tid, 4, "failed to schedule monitor check run, monitor", monitor.ID, err)
				return
			}
			cache.SetMonitorToRun(monitor.ID, *monitor.NextCheckAt)
		}
	})
	return nil
}
