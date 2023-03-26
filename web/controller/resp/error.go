package resp

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
)

const (
	ErrMessageValidationFailed = "Validation Failed"
	ErrMessageUnauthorized     = "Unauthorized"
)

const (
	ErrFailedToCreateGuestUser    = "failed to create guest user"
	ErrFailedToCreateOrganization = "failed to create organization"
	ErrGuestUserRateLimited       = "guest user rate limited"
	ErrGuestUserNotFound          = "guest user not found"
	ErrPlanNotFound               = "plan not found"
	ErrRoleNotFound               = "role not found"
	ErrCreatingNewUser            = "failed to create new user"
	ErrUpdatingUser               = "failed to update user"
	ErrDeletingGuestUser          = "failed to delete guest user"
	ErrGuestUserLoginFailed       = "failed to login guest user"
	ErrMonitorCreateFailed        = "failed to create monitor"
	ErrMalformedJWT               = "missing or malformed JWT"
)

var (
	ErrUsernameCannotBeEmpty        = errors.New("username cannot be empty")
	ErrPasswordCannotBeEmpty        = errors.New("password cannot be empty")
	ErrInvalidInterval              = errors.New("invalid interval")
	ErrInvalidBodyFormat            = errors.New("invalid body format")
	ErrInvalidAlarmReminderInterval = errors.New("invalid alarm reminder interval")
	ErrMaxBodySizeExceeded          = errors.New("max body size exceeded")
	ErrMaxTimeoutExceeded           = errors.New("max timeout exceeded")
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

func SendError(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
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

func ServeUnauthorizedError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": ErrMessageUnauthorized})
}
