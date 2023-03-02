package controller

import (
	"fmt"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
)

type UserController struct {
	userDomain *domain.UserDomain
}

func NewUserController(userDomain *domain.UserDomain) *UserController {
	return &UserController{userDomain: userDomain}
}

type GuestUserBody struct {
	Email string `json:"email" validate:"required,email,min=6,max=32"`
}

func (u *UserController) CreateGuestUser(c *fiber.Ctx) error {
	body := new(GuestUserBody)
	tracingID := pkg.GetTracingID(c.Context())

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	createUser := func() error {
		log.Default.Print(tracingID, 0, "creating guest user", body.Email)
		user, err := u.userDomain.CreateGuest(body.Email)
		if err != nil {
			return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToCreateGuestUser, err)
		}
		return resp.ServeData(c, fiber.StatusCreated, user)
	}

	latestGuestUser, err := u.userDomain.GetLatestGuestUser(body.Email)
	if err != nil {
		log.Default.Print(tracingID, 1, "no previous guest user", body.Email)
		return createUser()
	}

	passed := times.Now().Sub(latestGuestUser.InsertedAt).Seconds()
	remaining := int(constant.GuestUserRateLimitInSeconds - passed)
	log.Default.Print(tracingID, 2, "previous guest user exists", latestGuestUser.Email, "remaining", remaining)
	if passed <= constant.GuestUserRateLimitInSeconds {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrGuestUserRateLimited,
			fmt.Errorf("remaining %d", remaining))
	}
	return createUser()
}
