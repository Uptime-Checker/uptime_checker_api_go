package cron

import (
	"context"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
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
}
