package web

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func Setup(ctx context.Context) {
	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: config.IsProd,
	})

	// Sentry
	infra.SetupSentry()

	// NewRelic
	newRelicApp, err := infra.SetupNewRelic()
	if err != nil {
		sentry.CaptureException(err)
	}

	// Middlewares
	setupMiddlewares(app, newRelicApp)

	// Roues
	SetupRoutes(ctx, app)

	if err := app.Listen(fmt.Sprintf(":%s", config.App.Port)); err != nil {
		panic(err)
	}
}

func setupMiddlewares(app *fiber.App, newRelicApp *newrelic.Application) {
	app.Use(cors.New())
	app.Use(requestid.New(requestid.Config{
		ContextKey: string(constant.TracingKey), // => Setting Tracing ID to the context
		Generator: func() string {
			return pkg.GetUniqueString()
		},
	}))
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic: true,
	}))
	app.Use(logger.New(logger.Config{
		TimeZone: constant.UTCTimeZone,
		Format:   "[${time}] ${locals:tracing} | ${status} | ${latency} | ${method} | ${path}\n",
	}))
	app.Use(recover.New(recover.Config{EnableStackTrace: !config.IsProd}))

	if config.IsProd {
		app.Use(compress.New())
		app.Use(fibernewrelic.New(fibernewrelic.Config{
			Application: newRelicApp,
		}))
		app.Use(limiter.New(limiter.Config{
			Max: constant.MaxRequestPerMinute,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.Get(constant.OriginalIPHeader)
			},
		}))
	} else {
		app.Get("/metrics", monitor.New(monitor.Config{Title: "Metrics"}))
	}
}
