package domain

import (
	"context"
	"database/sql"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type CheckDomain struct{}

func NewCheckDomain() *CheckDomain {
	return &CheckDomain{}
}

func (c *CheckDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	monitorType resource.MonitorType,
) (*model.Monitor, error) {
	if !monitorType.Valid() {
		return nil, constant.ErrInvalidMonitorType
	}
	monitorTypeValue := int32(monitorType)

	if monitor.BodyFormat != nil {
		format := resource.MonitorBodyFormat(*monitor.BodyFormat)
		if !format.Valid() {
			return nil, constant.ErrInvalidMonitorBodyFormat
		}
	}

	monitor.Type = &monitorTypeValue
	insertStmt := Monitor.INSERT(Monitor.MutableColumns.Except(Monitor.InsertedAt, Monitor.UpdatedAt)).
		MODEL(monitor).
		RETURNING(Monitor.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, monitor)
	return monitor, err
}
