package controller

import (
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
	"github.com/Uptime-Checker/uptime_checker_api_go/task/client"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type MonitorController struct {
	monitorDomain       *domain.MonitorDomain
	regionDomain        *domain.RegionDomain
	monitorRegionDomain *domain.MonitorRegionDomain

	monitorService   *service.MonitorService
	assertionService *service.AssertionService

	dog *watchdog.WatchDog
}

func NewMonitorController(
	dog *watchdog.WatchDog,
	monitorDomain *domain.MonitorDomain,
	regionDomain *domain.RegionDomain,
	monitorRegionDomain *domain.MonitorRegionDomain,
	monitorService *service.MonitorService,
	assertionService *service.AssertionService,
) *MonitorController {
	return &MonitorController{
		dog:                 dog,
		monitorDomain:       monitorDomain,
		regionDomain:        regionDomain,
		monitorRegionDomain: monitorRegionDomain,
		assertionService:    assertionService,
		monitorService:      monitorService,
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
	// Body
	if body.Body != nil && len(*body.Body) > constant.MaxMonitorBodySizeInBytes {
		// Max 1KB
		return resp.ErrMaxBodySizeExceeded
	}

	// Timeout = max interval/2
	maxTimeoutRelativeToInterval := body.Interval / 2
	if body.Timeout > maxTimeoutRelativeToInterval {
		return resp.ErrMaxTimeoutExceeded
	}

	// Assertion
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
	// Validate the status code

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

	count, err := m.monitorDomain.Count(ctx, user.OrganizationID)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}

	lgr.Print(tracingID, 1, "monitor count", count, "org", user.Organization.Slug)
	if err := gandalf.CanCreateMonitor(user, int32(count), body.Interval); err != nil {
		return resp.SendError(c, fiber.StatusBadRequest, err)
	}

	var monitor *model.Monitor
	var assertions []*model.Assertion

	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
		lgr.Print(tracingID, 2, "creating monitor", body.Method, body.URL)
		monitor, err = m.monitorService.Create(ctx, tx, user.ID, user.OrganizationID, body.Name, body.URL, body.Method,
			body.Body, body.Username, body.Password, body.BodyFormat, body.Interval, body.Timeout,
			body.AlarmReminderInterval, body.AlarmReminderCount, body.CheckSSL,
			body.FollowRedirect, body.GlobalAlarmSettings, body.Headers)
		if err != nil {
			return err
		}

		// Insert the assertions
		for _, assertion := range body.Assertions {
			ass, err := m.assertionService.Create(ctx, tx, monitor.ID, assertion.Source, assertion.Property,
				assertion.Comparison, assertion.Value)
			if err != nil {
				return err
			}
			lgr.Print(tracingID, 3, "assertion created", resource.AssertionSource(ass.Source).String(), *ass.Value)
			assertions = append(assertions, ass)
		}
		return nil
	}); err != nil {
		lgr.Error(tracingID, 4, "failed to create monitor", err.Error())
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrMonitorCreateFailed, err)
	}

	// Start the monitor asynchronously
	if err := client.StartMonitorAsync(ctx, monitor.ID); err != nil {
		return resp.SendError(c, fiber.StatusInternalServerError, err)
	}

	return resp.ServeData(c, fiber.StatusOK, monitor)
}

func (m *MonitorController) ListMonitors(c *fiber.Ctx) error {
	user := middlelayer.GetUser(c)

	monitors, err := m.monitorDomain.List(c.Context(), user.OrganizationID, 5)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	return resp.ServeData(c, fiber.StatusOK, monitors)
}

type MonitorStartBody struct {
	MonitorID int64 `json:"monitorID" validate:"required"`
	On        bool  `json:"on"        validate:"required"`
}

func (m *MonitorController) Start(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(MonitorStartBody)
	if err := c.BodyParser(body); err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	if err := resp.Validate.Struct(body); err != nil {
		return resp.ServeValidationError(c, err)
	}

	if config.Region == nil {
		region, err := m.regionDomain.Get(ctx, config.App.FlyRegion)
		if err != nil {
			return err
		}
		config.Region = region
	}

	if body.On {
		monitorRegion, err := m.monitorRegionDomain.GetMonitorRegion(ctx, body.MonitorID, config.Region.ID)
		if err != nil {
			return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrMonitorNotFound, err)
		}
		monitorRegionWithAssertions, err := m.monitorRegionDomain.GetWithAllAssoc(ctx, monitorRegion.ID)
		if err != nil {
			return resp.ServeInternalServerError(c, err)
		}
		monitorWithAssertions := monitorRegionWithAssertions.Monitor
		assertions := monitorWithAssertions.Assertions
		m.dog.Start(ctx, monitorWithAssertions.Monitor, monitorRegionWithAssertions.Region, assertions)
	} else {
		if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
			_, err := m.monitorDomain.UpdateOn(ctx, tx, body.MonitorID, body.On, nil)
			return err
		}); err != nil {
			return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrMonitorCreateFailed, err)
		}
	}

	return resp.ServeNoContent(c, fiber.StatusNoContent)
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
	bodyFormat := resource.MonitorBodyFormatNoBody
	if body.BodyFormat != nil {
		resourceBodyFormat := resource.MonitorBodyFormat(*body.BodyFormat)
		bodyFormat = resourceBodyFormat
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
