package resp

import (
	"github.com/gofiber/fiber/v2"
)

type JSONResponse struct {
	Data any `json:"data,omitempty"`
}

func ServeData(c *fiber.Ctx, status int, data any) error {
	re := &JSONResponse{
		Data: data,
	}
	return c.Status(status).JSON(re)
}
