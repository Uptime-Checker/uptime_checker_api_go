package controller

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/gandalf"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/watchdog"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type MonitorController struct {
	monitorDomain *domain.MonitorDomain

	monitorService   *service.MonitorService
	assertionService *service.AssertionService

	dog *watchdog.WatchDog
}

func NewMonitorController(
	monitorDomain *domain.MonitorDomain,
	monitorService *service.MonitorService,
	assertionService *service.AssertionService,
	dog *watchdog.WatchDog,
) *MonitorController {
	return &MonitorController{
		monitorDomain:    monitorDomain,
		assertionService: assertionService,
		monitorService:   monitorService,
		dog:              dog,
	}
}

type AssertionBody struct {
	Source     int32   `json:"source"     validate:"required"`
	Property   *string `json:"property"`
	Comparison int32   `json:"comparison" validate:"required"`
	Value      string  `json:"value"      validate:"required"`
}

type MonitorBody struct {
	Name     string `json:"name"     validate:"required"`
	URL      string `json:"url"      validate:"required,url"`
	Method   string `json:"method"   validate:"required"`
	Interval int32  `json:"interval" validate:"required,min=10,max=86400"`
	Timeout  int32  `json:"timeout"  validate:"required,min=1,max=30"`

	Body       *string `json:"body"       validate:"required_with=bodyFormat"`
	BodyFormat *int32  `json:"bodyFormat" validate:"required_with=Body"`

	Headers *map[string]string `json:"headers"`

	Username *string `json:"username" validate:"required_with=password"`
	Password *string `json:"password" validate:"required_with=username"`

	GlobalAlarmSettings   bool  `json:"globalAlarmSettings"   validate:"required"`
	AlarmReminderInterval int32 `json:"alarmReminderInterval" validate:"required,min=600,max=3600"`
	AlarmReminderCount    int32 `json:"alarmReminderCount"    validate:"required,min=0,max=30"`

	CheckSSL       bool `json:"checkSSL"       validate:"required"`
	FollowRedirect bool `json:"followRedirect" validate:"required"`

	Assertions []AssertionBody `json:"assertions" validate:"required"`
}

func (m *MonitorController) validateMonitorBody(body *MonitorBody) error {
	// body
	if body.Body != nil && len(*body.Body) > constant.MaxMonitorBodySizeInBytes {
		// max 1KB
		return resp.ErrMaxBodySizeExceeded
	}

	// timeout = max interval/2
	maxTimeoutRelativeToInterval := body.Interval / 2
	if body.Timeout > maxTimeoutRelativeToInterval {
		return resp.ErrMaxTimeoutExceeded
	}

	// assertion
	statusCodeAssertionExists := false
	for _, assertion := range body.Assertions {
		if resource.AssertionSource(assertion.Source) == resource.AssertionSourceHeaders {
			if assertion.Property == nil {
				return resp.ErrHeaderKeyNeeded
			}
		}
		if resource.AssertionSource(assertion.Source) == resource.AssertionSourceStatusCode {
			statusCodeAssertionExists = true
		}
	}
	if !statusCodeAssertionExists {
		return resp.ErrStatusCodeAssertionRequired
	}
	// validate the status code

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
	if err := m.validateMonitorBody(body); err != nil {
		return resp.SendError(c, fiber.StatusUnprocessableEntity, err)
	}

	count, err := m.monitorDomain.Count(ctx, *user.OrganizationID)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	lgr.Print(tracingID, 1, "monitor count", count, "org", user.Organization.Slug)
	if err := gandalf.CanCreateMonitor(user, int32(count), body.Interval); err != nil {
		return resp.SendError(c, fiber.StatusBadRequest, err)
	}

	var monitor *model.Monitor
	if err := infra.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		lgr.Print(tracingID, 2, "creating monitor", body.Method, body.URL)
		monitor, err = m.monitorService.Create(ctx, tx, user.ID, *user.OrganizationID, body.Name, body.URL, body.Method,
			body.Body, body.Username, body.Password, body.BodyFormat, body.Interval, body.Timeout,
			body.AlarmReminderInterval, body.AlarmReminderCount, body.CheckSSL,
			body.FollowRedirect, body.GlobalAlarmSettings, body.Headers)
		if err != nil {
			return err
		}

		// insert the assertions
		for _, assertion := range body.Assertions {
			ass, err := m.assertionService.Create(ctx, tx, monitor.ID, assertion.Source, assertion.Property,
				assertion.Comparison, assertion.Value)
			if err != nil {
				return err
			}
			lgr.Print(tracingID, 3, "assertion created", resource.AssertionSource(*ass.Source).String(), *ass.Value)
		}
		return nil
	}); err != nil {
		lgr.Error(tracingID, 4, "failed to create monitor", err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrMonitorCreateFailed, err)
	}

	// start the monitor asynchronously
	go m.dog.Start(ctx, monitor, config.Region)

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

func (m *MonitorController) DryRun(c *fiber.Ctx) error {
	ctx := c.Context()
	body := new(MonitorBody)
	tracingID := pkg.GetTracingID(ctx)

	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}
	if err := m.validateMonitorBody(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	lgr.Print(tracingID, 1, "dry running", body.Method, body.URL)
	var bodyFormat *resource.MonitorBodyFormat
	if body.BodyFormat != nil {
		resourceBodyFormat := resource.MonitorBodyFormat(*body.BodyFormat)
		bodyFormat = &resourceBodyFormat
	}
	hitResponse, hitError := m.dog.Hit(ctx, body.URL, body.Method, body.Body, body.Username, body.Password, bodyFormat,
		body.Headers, body.Timeout, body.FollowRedirect)

	if hitResponse == nil && hitError != nil {
		err := fmt.Errorf("%s - %s", hitError.Type.String(), hitError.Text)
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrDryRunFailed, err)
	}

	for _, assertion := range body.Assertions {
		if pass := m.dog.Assert(
			assertion.Source, assertion.Property, assertion.Comparison, assertion.Value, *hitResponse,
		); !pass {
			err := fmt.Errorf("%s - value mismatch", resource.AssertionSource(assertion.Source).String())
			return resp.ServeDryRunError(c, fiber.StatusBadRequest, hitResponse, err)
		}
	}

	return resp.ServeData(c, fiber.StatusOK, hitResponse)
}
