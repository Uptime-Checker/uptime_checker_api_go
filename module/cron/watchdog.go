package cron

import (
	"context"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sourcegraph/conc/iter"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
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

	now := times.Now()
	prev := now.Add(time.Second * time.Duration(-watchDogCheckCronFromAndToInSeconds))
	later := now.Add(time.Second * time.Duration(+watchDogCheckCronFromAndToInSeconds))
	monitors, err := c.monitorDomain.ListMonitorsToRun(ctx, prev, later)
	if err != nil {
		return err
	}
	lgr.Print(tid, 2, "number of monitors to run:", len(monitors), "from", times.Format(prev), "to",
		times.Format(later), "region:", config.Region.Key)

	iterator := iter.Iterator[model.Monitor]{
		MaxGoroutines: watchDogCheckMaxGoroutine,
	}
	iterator.ForEach(monitors, func(monitor *model.Monitor) {
		cachedRegionID := cache.GetMonitorToRun(monitor.ID)
		if cachedRegionID == nil {
			// schedule the monitor
			currentMonitorRegion, err := c.monitorRegionDomain.GetMonitorRegion(ctx, monitor.ID, config.Region.ID)
			if err != nil {
				lgr.Error(tid, 3, "failed to get current monitor region", err)
				sentry.CaptureException(err)
				return
			}
			oldestCheckedMonitorRegion, err := c.monitorRegionDomain.GetOldestChecked(ctx, monitor.ID)
			if err != nil {
				lgr.Error(tid, 4, "failed to get the oldest monitor region", err)
				sentry.CaptureException(err)
				return
			}
			if currentMonitorRegion.ID == oldestCheckedMonitorRegion.ID {
				if err := client.RunCheckAsync(ctx, currentMonitorRegion.ID, *monitor.NextCheckAt); err != nil {
					lgr.Error(tid, 5, "failed to schedule monitor check run, monitor", monitor.ID, err)
					sentry.CaptureException(err)
					return
				}
				cache.SetMonitorToRun(monitor.ID, config.Region.ID)
			}
		}
	})
	return nil
}

func (c *Cron) stopTheDog(ctx context.Context) {
	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, 1, "check if we need to stopg the dog")
	watchDogValue := c.propertyService.Get(ctx, resource.PropertyKeyWatchDog)
	if watchDogValue == nil {
		runWatchDog, err := strconv.ParseBool(*watchDogValue)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		if !runWatchDog {
			s.Stop()
		}
	}
}
