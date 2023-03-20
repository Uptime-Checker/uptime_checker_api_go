package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type MonitorService struct {
	monitorDomain *domain.MonitorDomain
}

func NewMonitorService(monitorDomain *domain.MonitorDomain) *MonitorService {
	return &MonitorService{monitorDomain: monitorDomain}
}

func (m *MonitorService) Create(
	ctx context.Context,
	tx *sql.Tx,
	userID, organizationID int64,
	name, url, username, password string,
	method, interval, alarmReminderInterval, alarmReminderCount int32,
	checkSSL, followRedirect, globalAlarmSettings bool,
) (*model.Monitor, error) {
	head, getHeadErr := m.monitorDomain.GetHead(ctx, organizationID)

	now := times.Now()
	monitor := &model.Monitor{
		Name:                  name,
		URL:                   url,
		Method:                &method,
		Interval:              &interval,
		Username:              &username,
		Password:              &password,
		GlobalAlarmSettings:   &globalAlarmSettings,
		AlarmReminderInterval: &alarmReminderInterval,
		AlarmReminderCount:    &alarmReminderCount,
		CheckSsl:              &checkSSL,
		FollowRedirects:       &followRedirect,
		CreatedBy:             &userID,
		UpdatedBy:             &userID,
		PrevID:                nil,
		OrganizationID:        &organizationID,
	}

	monitor, err := m.monitorDomain.Create(ctx, tx, monitor, resource.MonitorTypeAPI)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
