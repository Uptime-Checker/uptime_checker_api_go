package controller

import (
	"github.com/gofiber/fiber/v2"
)

type GuestUserBody struct {
	Email string `validate:"required,email,min=6,max=32"`
}

func CreateGuestUser(c *fiber.Ctx) error {
	body := new(GuestUserBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(processValidationError(err))
	}
	return c.SendString("about")
}
