package web

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func Setup(ctx context.Context, shutdown context.CancelFunc) {
	tracingID := pkg.GetTracingID(ctx)

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

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", config.App.Port)); err != nil {
			panic(err)
		}
	}()

	quitCh := initQuitCh()
	sig := <-quitCh // This blocks the main thread until an interrupt is received
	lgr.Print(tracingID, "received", fmt.Sprintf("(%s)", sig.String()), "| gracefully shutting down...")
	if err := app.ShutdownWithTimeout(constant.ServerShutdownTimeoutInSeconds * time.Second); err != nil {
		sentry.CaptureException(err)
	}
	cleanup(ctx, shutdown)
	lgr.Print(tracingID, "app was successfully shutdown")
}

func initQuitCh() chan os.Signal {
	sigCh := make(chan os.Signal, 1) // Create channel to signify a signal being sent
	signal.Notify(
		sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	) // When an interrupt or termination signal is sent, notify the channel

	return sigCh
}

func cleanup(ctx context.Context, shutdown context.CancelFunc) {
	tracingID := pkg.GetTracingID(ctx)
	lgr.Print(tracingID, "running cleanup tasks...")

	// Shutdown the workers
	shutdown()

	// Close the DB connection
	if err := infra.DB.Close(); err != nil {
		sentry.CaptureException(err)
	}

	// Sync the logs
	if config.IsProd {
		lgr.Sync()
	}

	// Sync sentry
	infra.SyncSentry()
}

func setupMiddlewares(app *fiber.App, newRelicApp *newrelic.Application) {
	app.Use(cors.New())
	app.Use(requestid.New(requestid.Config{
		ContextKey: constant.TracingKey, // => Setting Tracing ID to the context
		Generator:  pkg.GetUniqueString,
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
		app.Use(pprof.New())
		app.Get("/metrics", monitor.New(monitor.Config{Title: "Metrics"}))
	}
}
