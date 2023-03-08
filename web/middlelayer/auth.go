package middlelayer

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
)

// Protected protect routes
func Protected(auth *service.AuthService) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearerToken := c.Get(constant.AuthorizationHeader)
		if !strings.HasPrefix(bearerToken, constant.AuthScheme) {
			return resp.ServeUnauthorizedError(c)
		}
		token := strings.TrimSpace(strings.TrimPrefix(bearerToken, constant.AuthScheme))
		if pkg.IsEmpty(token) {
			return resp.ServeUnauthorizedError(c)
		}

		user, err := auth.GetUserByToken(c.Context(), token)
		if err != nil {
			return resp.ServeUnauthorizedError(c)
		}
		c.Context().SetUserValue(string(constant.UserKey), user)
		return c.Next()
	}
}

// GetUser returns user that's in context
func GetUser(c *fiber.Ctx) *domain.UserWithRoleAndSubscription {
	v := c.Context().Value(string(constant.UserKey))
	if v == nil {
		panic("middleware: GetUser called without calling auth middleware prior")
	}
	u, _ := v.(*domain.UserWithRoleAndSubscription)
	return u
}
