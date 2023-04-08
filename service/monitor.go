package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
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
	headers *map[string]string,
) (*model.Monitor, error) {
	head, getHeadErr := m.monitorDomain.GetHead(ctx, organizationID)
	monitorMethod := resource.GetMonitorHTTPMethod(method)

	contentType := resource.MonitorBodyFormatNoBody
	if bodyFormat != nil {
		contentType = resource.MonitorBodyFormat(*bodyFormat)
	}
	pendingStatus := resource.MonitorStatusPending

	monitor := &model.Monitor{
		Name:                  name,
		URL:                   url,
		Method:                &monitorMethod,
		Timeout:               timeout,
		Interval:              interval,
		Body:                  body,
		Username:              username,
		Password:              password,
		On:                    false,
		Muted:                 false,
		GlobalAlarmSettings:   globalAlarmSettings,
		AlarmReminderInterval: alarmReminderInterval,
		AlarmReminderCount:    alarmReminderCount,
		Status:                int32(pendingStatus),
		CheckSsl:              checkSSL,
		FollowRedirects:       followRedirect,
		CreatedBy:             &userID,
		UpdatedBy:             &userID,
		NextID:                nil,
		OrganizationID:        organizationID,
	}

	if headers != nil && len(*headers) > 0 {
		jsonHeaders, err := json.Marshal(*headers)
		if err != nil {
			return nil, err
		}
		monitor.Headers = pkg.StringPointer(string(jsonHeaders))
	}

	// Create a new monitor
	monitor, err := m.monitorDomain.Create(ctx, tx, monitor, resource.MonitorTypeAPI, contentType)
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

	monitorStatusChange := &model.MonitorStatusChange{MonitorID: monitor.ID}

	// Create a new monitor status change
	_, err = m.monitorStatusDomain.Create(ctx, tx, monitorStatusChange, pendingStatus)
	if err != nil {
		return nil, err
	}

	return monitor, nil
}

func (m *MonitorService) StartOn(
	ctx context.Context,
	tx *sql.Tx,
	monitor *model.Monitor,
) (*model.Monitor, error) {
	passingStatus := resource.MonitorStatusPassing
	createMonitorStatusChange := func() error {
		monitorStatusChange := &model.MonitorStatusChange{MonitorID: monitor.ID}
		_, err := m.monitorStatusDomain.Create(ctx, tx, monitorStatusChange, passingStatus)
		return err
	}

	latestMonitorStatus, err := m.monitorStatusDomain.GetLatest(ctx, monitor.ID)
	if err != nil {
		if err := createMonitorStatusChange(); err != nil {
			return nil, fmt.Errorf("failed to create monitor status, err: %w", err)
		}
	}
	if resource.MonitorStatus(latestMonitorStatus.Status) != passingStatus {
		if err := createMonitorStatusChange(); err != nil {
			return nil, fmt.Errorf("failed to create monitor status, err: %w", err)
		}
	}

	now := times.Now()
	nextCheckAt := now.Add(time.Second * time.Duration(monitor.Interval))
	return m.monitorDomain.UpdateOnStatusAndCheckedAt(ctx, tx, monitor.ID, true, passingStatus, &now, &nextCheckAt)
}

func (m *MonitorService) GetRequestContentType(
	bodyFormat resource.MonitorBodyFormat,
	headers *map[string]string,
) string {
	contentType := bodyFormat.String()

	if headers != nil && len(*headers) > 0 {
		for key, value := range *headers {
			if strings.EqualFold(key, constant.ContentTypeHeader) {
				contentType = value
			}
		}
	}

	return contentType
}
