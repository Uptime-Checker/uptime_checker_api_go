package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type MonitorService struct {
	monitorDomain       *domain.MonitorDomain
	monitorStatusDomain *domain.MonitorStatusDomain
}

func NewMonitorService(
	monitorDomain *domain.MonitorDomain,
	monitorStatusDomain *domain.MonitorStatusDomain,
) *MonitorService {
	return &MonitorService{monitorDomain: monitorDomain, monitorStatusDomain: monitorStatusDomain}
}

func (m *MonitorService) Create(
	ctx context.Context,
	tx *sql.Tx,
	userID, organizationID int64,
	name, url, method string,
	body, username, password *string,
	bodyFormat *int32,
	interval, timeout, alarmReminderInterval, alarmReminderCount int32,
	checkSSL, followRedirect, globalAlarmSettings bool,
	headers map[string]string,
) (*model.Monitor, error) {
	head, getHeadErr := m.monitorDomain.GetHead(ctx, organizationID)
	monitorMethod := resource.GetMonitorHttpMethod(method)
	status := int32(resource.MonitorStatusPending)

	monitor := &model.Monitor{
		Name:                  name,
		URL:                   url,
		Method:                &monitorMethod,
		Timeout:               &timeout,
		Interval:              &interval,
		Body:                  body,
		BodyFormat:            bodyFormat,
		Username:              username,
		Password:              password,
		On:                    pkg.BoolPointer(true),
		Muted:                 pkg.BoolPointer(false),
		GlobalAlarmSettings:   &globalAlarmSettings,
		AlarmReminderInterval: &alarmReminderInterval,
		AlarmReminderCount:    &alarmReminderCount,
		Status:                &status,
		CheckSsl:              &checkSSL,
		FollowRedirects:       &followRedirect,
		CreatedBy:             &userID,
		UpdatedBy:             &userID,
		NextID:                nil,
		OrganizationID:        &organizationID,
	}

	if len(headers) > 0 {
		jsonHeaders, err := json.Marshal(headers)
		if err != nil {
			return nil, err
		}
		stringHeaders := string(jsonHeaders)
		monitor.Headers = &stringHeaders
	}

	// Create a new monitor
	monitor, err := m.monitorDomain.Create(ctx, tx, monitor, resource.MonitorTypeAPI)
	if err != nil {
		return nil, err
	}

	// Newly created monitor becomes the head, and we update previous head's next to the new head
	if head != nil && getHeadErr == nil {
		_, err = m.monitorDomain.UpdateNext(ctx, tx, head.ID, monitor.ID)
		if err != nil {
			return nil, err
		}
	}

	monitorStatusChange := &model.MonitorStatusChange{MonitorID: &monitor.ID}

	// Create a new monitor status change
	_, err = m.monitorStatusDomain.Create(ctx, tx, monitorStatusChange, resource.MonitorStatusPending)
	if err != nil {
		return nil, err
	}

	return monitor, nil
}
