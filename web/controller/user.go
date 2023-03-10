package controller

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/jinzhu/copier"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
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
		lgr.Default.Print(tracingID, 0, "creating guest user", body.Email)

		code := pkg.GetUniqueString()
		user, err := u.userDomain.CreateGuest(c.Context(), body.Email, pkg.HashSha(code))
		if err != nil {
			return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToCreateGuestUser, err)
		}
		return resp.ServeData(c, fiber.StatusCreated, user)
	}

	latestGuestUser, err := u.userDomain.GetLatestGuestUser(c.Context(), body.Email)
	if err != nil {
		lgr.Default.Print(tracingID, 1, "no previous guest user", body.Email)
		return createUser()
	}

	passed := times.Now().Sub(latestGuestUser.InsertedAt).Minutes()
	remaining := int(constant.GuestUserRateLimitInMinutes - passed)
	lgr.Default.Print(tracingID, 2, "previous guest user exists", latestGuestUser.Email, "remaining", remaining)
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
	lgr.Default.Print(tracingID, 1, "found guest user", guestUser.ID, guestUser.Email, guestUser.ExpiresAt)

	user, userGetError := u.userDomain.GetUser(ctx, body.Email)
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if userGetError != nil {
			user, err = u.userService.CreateNewUserAndContact(ctx, tx, body.Email)
			if err != nil {
				return err
			}
			lgr.Default.Print(tracingID, 2, "created new user", user.ID, user.Email, "provider", *user.Provider)
		} else {
			user, err = u.userDomain.UpdateProvider(ctx, tx, user.ID, body.Email, resource.UserLoginProviderEmail)
			if err != nil {
				return err
			}
			lgr.Default.Print(tracingID, 3, "update user provider", user.ID, user.Email, user.Provider)
		}
		lgr.Default.Print(tracingID, 4, "deleting guest user", guestUser.ID, guestUser.Email)
		return u.userDomain.DeleteGuestUser(ctx, tx, guestUser.ID)
	}); err != nil {
		lgr.Default.Error(tracingID, 5, "failed to login guest user", err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrGuestUserLoginFailed, err)
	}

	token, err := u.authService.GenerateUserToken(user)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	respUser := resp.User{Token: token}
	if err := copier.Copy(&respUser, &user); err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	return resp.ServeData(c, fiber.StatusOK, respUser)
}

type ProviderLoginBody struct {
	Name        string `json:"name"        validate:"required"`
	Email       string `json:"email"       validate:"required,email,min=6,max=32"`
	Provider    int    `json:"provider"    validate:"required"`
	ProviderUID string `json:"providerUID" validate:"required"`
	Picture     string `json:"picture"     validate:"url"`
}

func (u *UserController) ProviderLogin(c *fiber.Ctx) error {
	ctx := c.Context()
	body := new(ProviderLoginBody)
	tracingID := pkg.GetTracingID(ctx)

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	var user *model.User
	var err error
	provider := resource.UserLoginProvider(body.Provider)

	user, err = u.userDomain.GetUser(ctx, body.Email)
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if err != nil {
			user, err = u.userService.CreateNewProviderUserAndContact(
				ctx, tx, body.Name, body.Email, body.Provider, body.ProviderUID, body.Picture,
			)
			if err != nil {
				return err
			}
			lgr.Default.Print(tracingID, 1, "created new user", user.ID, user.Email, "provider", provider.String())
		} else {
			user, err = u.userDomain.UpdateProvider(
				ctx, tx, user.ID, body.Email, resource.UserLoginProvider(body.Provider),
			)
			if err != nil {
				return err
			}
			lgr.Default.Print(tracingID, 2, "update user provider", user.ID, user.Email, "provider", provider.String())
		}
		return nil
	}); err != nil {
		lgr.Default.Error(tracingID, 3, "failed to login provider user", body.Email, provider.String(), err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrGuestUserLoginFailed, err)
	}
	token, err := u.authService.GenerateUserToken(user)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	respUser := resp.User{Token: token}
	if err := copier.Copy(&respUser, &user); err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	return resp.ServeData(c, fiber.StatusOK, respUser)
}

func (u *UserController) GetMe(c *fiber.Ctx) error {
	user := middlelayer.GetUser(c)
	return resp.ServeData(c, fiber.StatusOK, user)
}

type UserUpdateBody struct {
	Name string `json:"name" validate:"required,min=4,max=32"`
}

func (u *UserController) Update(c *fiber.Ctx) error {
	user := middlelayer.GetUser(c)
	body := new(UserUpdateBody)

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	updatedUser, err := u.userDomain.UpdateName(c.Context(), user.User.ID, body.Name)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrUpdatingUser, err)
	}
	return resp.ServeData(c, fiber.StatusOK, updatedUser)
}
