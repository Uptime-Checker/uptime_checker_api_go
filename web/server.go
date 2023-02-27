package web

import (
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
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func Setup() {
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
	SetupRoutes(app)

	if err := app.Listen(fmt.Sprintf(":%s", config.App.Port)); err != nil {
		panic("Server start failed")
	}
}

func setupMiddlewares(app *fiber.App, newRelicApp *newrelic.Application) {
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(requestid.New(requestid.Config{
		ContextKey: string(constant.TracingKey), // => Setting Tracing ID to the context
		Generator: func() string {
			return pkg.GetUniqueString()
		},
	}))
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic: true,
	}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return !config.IsProd // => Skip limiter in dev
		},
		Max: constant.MaxRequestPerMinute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get(constant.OriginalIPHeader)
		},
	}))
	app.Use(logger.New(logger.Config{
		TimeZone: constant.UTCTimeZone,
		Format:   "[${time}] ${locals:tracing} | ${status} | ${latency} | ${method} | ${path}\n",
	}))

	if config.IsProd {
		app.Use(fibernewrelic.New(fibernewrelic.Config{
			Application: newRelicApp,
		}))
	}
}
