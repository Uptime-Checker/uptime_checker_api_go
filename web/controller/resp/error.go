package resp

import (
	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
)

const (
	ErrMessageValidationFailed = "Validation Failed"
)

const (
	ErrFailedToCreateGuestUser = "failed to create guest user"
	ErrGuestUserRateLimited    = "guest user rate limited"
	ErrGuestUserNotFound       = "guest user not found"
)

var Validate = validator.New()

type ValidationError struct {
	Error  string      `json:"error"`
	Field  string      `json:"field"`
	Reason string      `json:"reason"`
	Value  interface{} `json:"value"`
}

func processValidationError(err error) ValidationError {
	validationErr := ValidationError{}
	for _, err := range err.(validator.ValidationErrors) {
		validationErr.Error = ErrMessageValidationFailed
		validationErr.Field = err.Field()
		validationErr.Reason = err.Tag()
		validationErr.Value = err.Value()

		// Just return one validation error at once
		break
	}
	return validationErr
}

func processError(message string, err error) map[string]interface{} {
	return fiber.Map{"message": message, "error": err.Error()}
}

func ServeError(c *fiber.Ctx, status int, message string, err error) error {
	return c.Status(status).JSON(processError(message, err))
}

func ServeInternalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}

func ServeValidationError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(processValidationError(err))
}
