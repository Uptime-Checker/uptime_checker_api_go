package controller

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type MonitorIntegrationController struct {
	monitorIntegrationService *service.MonitorIntegrationService
}

func NewMonitorIntegrationController(
	monitorIntegrationService *service.MonitorIntegrationService,
) *MonitorIntegrationController {
	return &MonitorIntegrationController{monitorIntegrationService: monitorIntegrationService}
}

type MonitorIntegrationBody struct {
	MonitorIntegrationType int32             `json:"type"   validate:"required"`
	Config                 map[string]string `json:"config" validate:"required"`
}

func (m *MonitorIntegrationController) Create(c *fiber.Ctx) error {
	ctx := c.Context()
	user := middlelayer.GetUser(c)
	body := new(MonitorIntegrationBody)
	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	integrationType := resource.MonitorIntegrationType(body.MonitorIntegrationType)
	if integrationType == resource.MonitorIntegrationTypeWebhook {
		url := body.Config["url"]
		if pkg.IsEmpty(url) {
			return resp.SendError(c, fiber.StatusUnprocessableEntity, resp.ErrWebhookURLRequired)
		}
	}

	var err error
	var monitorIntegration *model.MonitorIntegration
	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
		monitorIntegration, err = m.monitorIntegrationService.Create(
			ctx, tx, *user.OrganizationID, integrationType, body.Config,
		)
		return err
	}); err != nil {
		return resp.SendError(c, fiber.StatusBadRequest, err)
	}
	return resp.ServeData(c, fiber.StatusOK, monitorIntegration)
}
