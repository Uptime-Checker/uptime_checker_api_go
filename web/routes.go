package web

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

func SetupRoutes(app *fiber.App) {
	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// API V1
	v1 := app.Group("/v1")
	v1.Use(middlelayer.Header())
	v1.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Domain Registration
	userDomain := domain.NewUserDomain()

	// User router for auth and user account
	userRouter := v1.Group("/user")
	registerUserHandlers(userRouter, userDomain)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
}

func registerUserHandlers(router fiber.Router, userDomain *domain.UserDomain) {
	handler := controller.NewUserController(userDomain)
	router.Post("/guest", handler.CreateGuestUser)
}
