package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type MonitorStatusDomain struct{}

func NewMonitorStatusDomain() *MonitorStatusDomain {
	return &MonitorStatusDomain{}
}

func (m *MonitorStatusDomain) List(
	ctx context.Context,
	monitorID int64,
	limit int,
) ([]model.MonitorStatusChange, error) {
	stmt := SELECT(
		MonitorStatusChange.AllColumns,
	).FROM(MonitorStatusChange).
		WHERE(MonitorStatusChange.MonitorID.EQ(Int(monitorID))).
		ORDER_BY(MonitorStatusChange.InsertedAt.DESC()).LIMIT(int64(limit))

	var monitorStatuses []model.MonitorStatusChange
	err := stmt.QueryContext(ctx, infra.DB, &monitorStatuses)
	return monitorStatuses, err
}

func (m *MonitorStatusDomain) GetLatest(
	ctx context.Context,
	monitorID int64,
) (*model.MonitorStatusChange, error) {
	stmt := SELECT(
		MonitorStatusChange.AllColumns,
	).FROM(MonitorStatusChange).
		WHERE(MonitorStatusChange.MonitorID.EQ(Int(monitorID))).
		ORDER_BY(MonitorStatusChange.InsertedAt.DESC()).LIMIT(1)

	monitorStatus := &model.MonitorStatusChange{}
	err := stmt.QueryContext(ctx, infra.DB, monitorStatus)
	return monitorStatus, err
}

func (m *MonitorStatusDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	monitorStatusChange *model.MonitorStatusChange,
	monitorStatus resource.MonitorStatus,
) (*model.MonitorStatusChange, error) {
	if !monitorStatus.Valid() {
		return nil, constant.ErrInvalidMonitorStatus
	}

	monitorStatusChange.Status = int32(monitorStatus)
	insertStmt := MonitorStatusChange.INSERT(MonitorStatusChange.MutableColumns.
		Except(MonitorStatusChange.InsertedAt, MonitorStatusChange.UpdatedAt)).
		MODEL(monitorStatusChange).RETURNING(MonitorStatusChange.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, monitorStatusChange)
	return monitorStatusChange, err
}
