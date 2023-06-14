package domain

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type MonitorNotificationnDomain struct{}

func NewMonitorNotificationDomain() *MonitorNotificationnDomain {
	return &MonitorNotificationnDomain{}
}

func (m *MonitorNotificationnDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	monitorNotification *model.MonitorNotification,
	monitorNotificationType resource.MonitorNotificationType,
) (*model.MonitorNotification, error) {
	if !monitorNotificationType.Valid() {
		return nil, constant.ErrInvalidMonitorIntegrationType
	}
	monitorNotification.Type = int32(monitorNotificationType)
	insertStmt := MonitorNotification.INSERT(MonitorNotification.MutableColumns.
		Except(MonitorNotification.InsertedAt, MonitorNotification.UpdatedAt)).MODEL(monitorNotification).
		RETURNING(MonitorNotification.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, monitorNotification)
	return monitorNotification, err
}
