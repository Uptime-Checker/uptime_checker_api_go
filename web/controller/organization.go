package controller

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
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
	log.Default.Print(tracingID, 1, "found plan", plan.Name, "product", plan.Product.Name)

	role, err := o.organizationDomain.GetRoleByType(ctx, resource.RoleTypeSuperAdmin)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrRoleNotFound, err)
	}
	log.Default.Print(tracingID, 2, "to assign role", role.Name)

	var organization *model.Organization
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		organization, err = o.organizationService.Create(ctx, tx, body.Name, body.Slug, user.User.ID, role.ID)
		if err != nil {
			return err
		}
		log.Default.Print(tracingID, 3, "created organization", organization.Name, "slug", organization.Slug)

		updatedUser, err := o.userDomain.UpdateOrganizationAndRole(ctx, tx, user.User.ID, role.ID, organization.ID)
		if err != nil {
			return err
		}
		log.Default.Print(tracingID, 4, "updated user role", updatedUser.ID, "organization",
			organization.Slug, "role", role.Name)

		subscription, err := o.paymentService.CreateSubscription(ctx, tx, organization.ID, *plan)
		log.Default.Print(tracingID, 5, "subscription started", subscription.ID, "plan", plan.Name,
			"product", plan.Product.Name)
		return err
	}); err != nil {
		log.Default.Error(tracingID, 6, "failed to create organization", err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToCreateOrganization, err)
	}
	return resp.ServeData(c, fiber.StatusOK, organization)
}
