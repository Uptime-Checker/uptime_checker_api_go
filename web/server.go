package web

import (
	"fmt"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Setup() {
	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: prod(config.App.Release),
	})

	// Middlewares
	setupMiddlewares(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	if err := app.Listen(fmt.Sprintf(":%s", config.App.Port)); err != nil {
		panic("Server start failed")
	}
}

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(requestid.New())
	app.Use(limiter.New(limiter.Config{
		Max: constant.MaxRequestPerMinute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get(constant.OriginalIPHeader)
		},
	}))
}

func prod(release string) bool {
	return release == constant.ReleaseProd
}
