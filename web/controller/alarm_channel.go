package controller

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"

	"github.com/gofiber/fiber/v2"
)

type AlarmChannelController struct {
	alarmChannelDomain *domain.AlarmChannelDomain
}

func NewAlarmChannelController(
	alarmChannelDomain *domain.AlarmChannelDomain,
) *AlarmChannelController {
	return &AlarmChannelController{
		alarmChannelDomain: alarmChannelDomain,
	}
}

func (m *AlarmChannelController) List(c *fiber.Ctx) error {
	ctx := c.Context()
	user := middlelayer.GetUser(c)

	monitorID := int64(c.QueryInt("monitor_id"))
	if monitorID != 0 {
		alarmChannels, err := m.alarmChannelDomain.ListByMonitor(ctx, monitorID)
		if err != nil {
			return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToListAlarmChannels, err)
		}
		return resp.ServeData(c, fiber.StatusOK, alarmChannels)
	}

	alarmChannels, err := m.alarmChannelDomain.ListByOrganization(ctx, *user.OrganizationID)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrFailedToListAlarmChannels, err)
	}
	return resp.ServeData(c, fiber.StatusOK, alarmChannels)
}
