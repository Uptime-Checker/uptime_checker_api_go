package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: false, // go run app.go -prod
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err := app.Listen(":3000")
	if err != nil {
		panic("Server start failed")
	}
}
