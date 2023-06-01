package controller

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type MonitorIntegrationController struct {
	alarmChannelDomain        *domain.AlarmChannelDomain
	monitorIntegrationService *service.MonitorIntegrationService
}

func NewMonitorIntegrationController(
	alarmChannelDomain *domain.AlarmChannelDomain,
	monitorIntegrationService *service.MonitorIntegrationService,
) *MonitorIntegrationController {
	return &MonitorIntegrationController{
		alarmChannelDomain:        alarmChannelDomain,
		monitorIntegrationService: monitorIntegrationService,
	}
}

type MonitorIntegrationBody struct {
	MonitorIntegrationType int32          `json:"type"   validate:"required"`
	Config                 map[string]any `json:"config" validate:"required"`
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
		_, ok := body.Config["url"].(string)
		if !ok {
			return resp.SendError(c, fiber.StatusUnprocessableEntity, resp.ErrWebhookURLRequired)
		}
	} else if integrationType == resource.MonitorIntegrationTypeSlack {
		_, ok := body.Config["access_token"].(string)
		if !ok {
			return resp.SendError(c, fiber.StatusUnprocessableEntity, resp.ErrAccessTokenRequired)
		}
		_, ok = body.Config["incoming_webhook"].(map[string]any)
		if !ok {
			return resp.SendError(c, fiber.StatusUnprocessableEntity, resp.ErrIncomingWebhookRequired)
		}
	}

	var err error
	var monitorIntegration *model.MonitorIntegration
	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
		monitorIntegration, err = m.monitorIntegrationService.Create(
			ctx, tx, user.Organization, integrationType, body.Config,
		)
		if err != nil {
			return err
		}
		alarmChannel := &model.AlarmChannel{
			On:             true,
			OrganizationID: user.Organization.ID,
			IntegrationID:  &monitorIntegration.ID,
		}
		_, err = m.alarmChannelDomain.Create(ctx, tx, alarmChannel)
		return err
	}); err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToCreateIntegration, err)
	}
	return resp.ServeData(c, fiber.StatusOK, monitorIntegration)
}
