package controller

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/gandalf"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
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

	Body       string `json:"body"`
	BodyFormat string `json:"bodyFormat"`

	Headers map[string]string `mapstructure:"headers"`

	Username *string `json:"username"`
	Password *string `json:"password"`

	GlobalAlarmSettings   bool `json:"globalAlarmSettings"   validate:"required"`
	AlarmReminderInterval int  `json:"alarmReminderInterval" validate:"required"`
	AlarmReminderCount    int  `json:"alarmReminderCount"    validate:"required"`

	CheckSSL       bool `json:"checkSSL"       validate:"required"`
	FollowRedirect bool `json:"followRedirect" validate:"required"`
}

func (m *MonitorController) validateMonitorBody(body *MonitorBody) error {

	return nil
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

	lgr.Default.Print(tracingID, 1, "monitor count", count, "org", user.Organization.Slug)
	if err := gandalf.CanCreateMonitor(user, int32(count), int32(body.Interval)); err != nil {
		return resp.SendError(c, fiber.StatusBadRequest, err)
	}

	var monitor *model.Monitor
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		lgr.Default.Print(tracingID, 2, "creating monitor", body.Method, body.URL)
		monitor, err = m.monitorService.Create(ctx, tx, user.ID, *user.OrganizationID, body.Name, body.URL,
			body.Method, body.Username, body.Password, int32(body.Interval), int32(body.AlarmReminderInterval),
			int32(body.AlarmReminderCount), body.CheckSSL, body.FollowRedirect, body.GlobalAlarmSettings)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		lgr.Default.Error(tracingID, 3, "failed to create monitor", err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrMonitorCreateFailed, err)
	}

	return resp.ServeData(c, fiber.StatusOK, monitor)
}

func (m *MonitorController) ListMonitors(c *fiber.Ctx) error {
	user := middlelayer.GetUser(c)

	monitors, err := m.monitorDomain.List(c.Context(), *user.OrganizationID, 5)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	return resp.ServeData(c, fiber.StatusOK, monitors)
}
