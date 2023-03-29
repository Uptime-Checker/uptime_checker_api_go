package infra

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func SetupSentry() {
	options := sentry.ClientOptions{
		Dsn:              config.App.SentryDSN,
		Environment:      config.App.Release,
		Release:          config.App.Version,
		TracesSampleRate: 0.2,
		AttachStacktrace: config.IsProd,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
					tracingID := pkg.GetTracingID(c.Context())
					event.Extra = map[string]any{string(constant.TracingKey): tracingID}
				}
			}
			return event
		},
	}
	options.EnableTracing = config.IsProd
	options.Debug = !config.IsProd

	err := sentry.Init(options)
	if err != nil {
		lgr.Errorf("sentry.Init: %s", err)
	}
}

func SyncSentry() {
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	sentry.Flush(1 * time.Second)
}
