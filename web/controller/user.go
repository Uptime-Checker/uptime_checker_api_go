package controller

import (
	"github.com/gofiber/fiber/v2"
)

type GuestUserBody struct {
	Email string `validate:"required,email,min=6,max=32"`
}

func CreateGuestUser(c *fiber.Ctx) error {
	return c.SendString("about")
}
