package web

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func Setup() {
	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: config.IsProd,
	})

	// Middlewares
	setupMiddlewares(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	if err := app.Listen(fmt.Sprintf(":%s", config.App.Port)); err != nil {
		panic("Server start failed")
	}
}

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(requestid.New(requestid.Config{
		ContextKey: string(constant.TracingKey), // => Setting Tracing ID to the context
		Generator: func() string {
			return pkg.GetUniqueString()
		},
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
}
