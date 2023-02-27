package middlelayer

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
)

func Header() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		headerKey := c.Get(constant.APIKeyHeader)
		if headerKey != config.App.APIKey {
			return fiber.ErrPreconditionFailed
		}

		return c.Next()
	}
}
