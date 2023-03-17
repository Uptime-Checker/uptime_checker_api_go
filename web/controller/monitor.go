package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/gandalf"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type MonitorController struct {
	monitorDomain  *domain.MonitorDomain
	monitorService *service.MonitorService
}

func NewMonitorController(
	monitorDomain *domain.MonitorDomain,
	monitorService *service.MonitorService,
) *MonitorController {
	return &MonitorController{monitorDomain: monitorDomain, monitorService: monitorService}
}

type MonitorBody struct {
	Name     string `json:"name"     validate:"required"`
	URL      string `json:"url"      validate:"required,url"`
	Method   string `json:"method"   validate:"required"`
	Interval int    `json:"interval" validate:"required"`

	Body       string            `json:"body"`
	BodyFormat string            `json:"bodyFormat"`
	Headers    map[string]string `mapstructure:"headers"`

	Username string `json:"username"`
	Password string `json:"password"`

	CheckSSL       bool `json:"checkSSL"`
	FollowRedirect bool `json:"followRedirect"`
}

func (m *MonitorController) Create(c *fiber.Ctx) error {
	ctx := c.Context()
	body := new(MonitorBody)
	tracingID := pkg.GetTracingID(ctx)
	user := middlelayer.GetUser(c)

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	count, err := m.monitorDomain.Count(ctx, *user.OrganizationID)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	if err := gandalf.CanCreateMonitor(user, int32(count), int32(body.Interval)); err != nil {
		return resp.SendError(c, fiber.StatusBadRequest, err)
	}

	return c.SendString(tracingID)
}

func (m *MonitorController) ListMonitors(c *fiber.Ctx) error {
	return nil
}
