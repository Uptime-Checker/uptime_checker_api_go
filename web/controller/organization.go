package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
)

type OrganizationController struct {
	paymentDomain *domain.PaymentDomain
	userService   *service.UserService
}

func NewOrganizationController(
	paymentDomain *domain.PaymentDomain,
	userService *service.UserService,
) *OrganizationController {
	return &OrganizationController{paymentDomain: paymentDomain, userService: userService}
}

type CreateOrganizationBody struct {
	Name   string `json:"name"   validate:"required,min=4,max=32"`
	Slug   string `json:"slug"   validate:"required,min=4,max=32"`
	PlanID int64  `json:"planID" validate:"required"`
}

func (o *OrganizationController) CreateOrganization(c *fiber.Ctx) error {
	ctx := c.Context()
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
	log.Default.Print(tracingID, 1, "found plan", plan.Name, plan.Product.Name)
	return resp.ServeData(c, fiber.StatusOK, plan)
}
