package main

import (
	"fmt"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		panic("Config load failed")
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: prod(config.Release),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		port := viper.GetString("PORT")
		return c.SendString(port)
	})

	if err := app.Listen(fmt.Sprintf(":%s", config.Port)); err != nil {
		panic("Server start failed")
	}
}

func prod(release string) bool {
	return release == constant.ReleaseProd
}
