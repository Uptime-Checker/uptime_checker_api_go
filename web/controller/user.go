package controller

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
)

type UserController struct {
	userDomain  *domain.UserDomain
	authService *service.AuthService
	userService *service.UserService
}

func NewUserController(
	userDomain *domain.UserDomain,
	authService *service.AuthService,
	userService *service.UserService,
) *UserController {
	return &UserController{userDomain: userDomain, authService: authService, userService: userService}
}

type CreateGuestUserBody struct {
	Email string `json:"email" validate:"required,email,min=6,max=32"`
}

func (u *UserController) CreateGuestUser(c *fiber.Ctx) error {
	body := new(CreateGuestUserBody)
	tracingID := pkg.GetTracingID(c.Context())

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	createUser := func() error {
		log.Default.Print(tracingID, 0, "creating guest user", body.Email)

		code := pkg.GetUniqueString()
		user, err := u.userDomain.CreateGuest(c.Context(), body.Email, pkg.HashSha(code))
		if err != nil {
			return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToCreateGuestUser, err)
		}
		return resp.ServeData(c, fiber.StatusCreated, user)
	}

	latestGuestUser, err := u.userDomain.GetLatestGuestUser(c.Context(), body.Email)
	if err != nil {
		log.Default.Print(tracingID, 1, "no previous guest user", body.Email)
		return createUser()
	}

	passed := times.Now().Sub(latestGuestUser.InsertedAt).Minutes()
	remaining := int(constant.GuestUserRateLimitInMinutes - passed)
	log.Default.Print(tracingID, 2, "previous guest user exists", latestGuestUser.Email, "remaining", remaining)
	if passed <= constant.GuestUserRateLimitInMinutes {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrGuestUserRateLimited,
			fmt.Errorf("remaining %d", remaining))
	}
	return createUser()
}

type GuestUserLoginBody struct {
	Email string `json:"email" validate:"required,email,min=6,max=32"`
	Code  string `json:"code"  validate:"required,min=6,max=100"`
}

func (u *UserController) GuestUserLogin(c *fiber.Ctx) error {
	ctx := c.Context()
	body := new(GuestUserLoginBody)
	tracingID := pkg.GetTracingID(ctx)

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	guestUser, err := u.authService.VerifyGuestUser(ctx, body.Email, body.Code)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrGuestUserNotFound, err)
	}
	log.Default.Print(tracingID, 1, "found guest user", guestUser.ID, guestUser.Email, guestUser.ExpiresAt)

	user, err := u.userDomain.GetUser(ctx, body.Email)
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err != nil {
			user, err = u.userService.CreateNewUserAndContact(ctx, tx, body.Email)
			if err != nil {
				return err
			}
			log.Default.Print(tracingID, 2, "created new user", user.ID, user.Email, "provider", *user.Provider)
		} else {
			user, err = u.userDomain.UpdateProvider(ctx, tx, body.Email, resource.UserLoginProviderEmail)
			if err != nil {
				return err
			}
			log.Default.Print(tracingID, 3, "update user provider", user.ID, user.Email, user.Provider)
		}
		log.Default.Print(tracingID, 4, "deleting guest user", guestUser.ID, guestUser.Email)
		return u.userDomain.DeleteGuestUser(ctx, tx, guestUser.ID)
	}); err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrGuestUserLoginFailed, err)
	}

	return resp.ServeData(c, fiber.StatusOK, user)
}
