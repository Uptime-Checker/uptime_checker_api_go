package controller

import (
	"context"
	"database/sql"

	"github.com/Joker666/future"
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type OrganizationController struct {
	userDomain         *domain.UserDomain
	paymentDomain      *domain.PaymentDomain
	organizationDomain *domain.OrganizationDomain

	paymentService      *service.PaymentService
	organizationService *service.OrganizationService
}

func NewOrganizationController(
	userDomain *domain.UserDomain,
	paymentDomain *domain.PaymentDomain,
	organizationDomain *domain.OrganizationDomain,
	paymentService *service.PaymentService,
	organizationService *service.OrganizationService,
) *OrganizationController {
	return &OrganizationController{
		userDomain:          userDomain,
		paymentDomain:       paymentDomain,
		organizationDomain:  organizationDomain,
		paymentService:      paymentService,
		organizationService: organizationService,
	}
}

type CreateOrganizationBody struct {
	Name   string `json:"name"   validate:"required,min=4,max=32"`
	Slug   string `json:"slug"   validate:"required,min=4,max=32"`
	PlanID int64  `json:"planID" validate:"required"`
}

func (o *OrganizationController) CreateOrganization(c *fiber.Ctx) error {
	ctx := c.Context()
	user := middlelayer.GetUser(c)
	body := new(CreateOrganizationBody)
	tracingID := pkg.GetTracingID(ctx)

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	plan, err := o.paymentDomain.GetPlanWithProduct(ctx, body.PlanID)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrPlanNotFound, err)
	}
	lgr.Print(tracingID, 1, "found plan", plan.Name, "product", plan.Product.Name)

	role, err := o.organizationDomain.GetRoleByType(ctx, resource.RoleTypeSuperAdmin)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrRoleNotFound, err)
	}
	lgr.Print(tracingID, 2, "to assign role", role.Name)

	var organization *model.Organization
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		organization, err = o.organizationDomain.Create(ctx, tx, body.Name, body.Slug)
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 3, "created organization", organization.Name, "slug", organization.Slug)

		organizationUserAsync := future.New(func() (*model.OrganizationUser, error) {
			return o.organizationDomain.CreateOrganizationUser(ctx, tx, organization.ID, user.ID, *user.RoleID)
		})
		updateOrganizationAndRoleAsync := future.New(func() (*model.User, error) {
			return o.userDomain.UpdateOrganizationAndRole(ctx, tx, user.User.ID, role.ID, organization.ID)
		})
		subscriptionAsync := future.New(func() (*model.Subscription, error) {
			return o.paymentService.CreateSubscription(ctx, tx, organization.ID, *plan)
		})
		alarmPolicyAsync := future.New(func() (*model.AlarmPolicy, error) {
			return o.organizationService.CreateOrganizationAlarmPolicy(ctx, tx, organization.ID)
		})

		organizationUser, err := organizationUserAsync.Await()
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 4, "organization user created", organizationUser.ID)

		updatedUser, err := updateOrganizationAndRoleAsync.Await()
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 5, "updated user role", updatedUser.ID, "organization", organization.Slug,
			"role", role.Name)

		subscription, err := subscriptionAsync.Await()
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 6, "subscription started", subscription.ID, "plan", plan.Name,
			"product", plan.Product.Name)

		alarmPolicy, err := alarmPolicyAsync.Await()
		if err != nil {
			return err
		}
		lgr.Print(tracingID, 7, "organization alarm policy created", alarmPolicy.Reason, alarmPolicy.Threshold)

		return nil
	}); err != nil {
		lgr.Error(tracingID, 8, "failed to create organization", err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToCreateOrganization, err)
	}

	cache.DeleteUserWithRoleAndSubscription(user.ID)
	return resp.ServeData(c, fiber.StatusOK, organization)
}

func (o *OrganizationController) ListOrganizationsOfUser(c *fiber.Ctx) error {
	user := middlelayer.GetUser(c)
	organizationUserRoles, err := o.organizationDomain.ListOrganizationsOfUser(c.Context(), user.ID)
	if err != nil {
		return resp.SendError(c, fiber.StatusBadRequest, err)
	}
	return resp.ServeData(c, fiber.StatusOK, organizationUserRoles)
}
