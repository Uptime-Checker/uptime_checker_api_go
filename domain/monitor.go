package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type MonitorDomain struct{}

func NewMonitorDomain() *MonitorDomain {
	return &MonitorDomain{}
}

func (m *MonitorDomain) Count(ctx context.Context, organizationID int64) (int, error) {
	stmt := SELECT(COUNT(Monitor.ID)).FROM(Monitor).WHERE(Monitor.OrganizationID.EQ(Int(organizationID)))

	var dest struct {
		count int
	}
	err := stmt.QueryContext(ctx, infra.DB, &dest)
	return dest.count, err
}

func (m *MonitorDomain) List(ctx context.Context, organizationID int64, limit int) ([]model.Monitor, error) {
	stmt := m.listRecursively(organizationID, limit)

	var monitors []model.Monitor
	err := stmt.QueryContext(ctx, infra.DB, &monitors)
	return monitors, err
}

func (m *MonitorDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	monitor *model.Monitor,
	monitorType resource.MonitorType,
) (*model.Monitor, error) {
	if !monitorType.Valid() {
		return nil, constant.ErrInvalidMonitorType
	}
	monitorTypeValue := int32(monitorType)

	monitor.Type = &monitorTypeValue
	insertStmt := Monitor.INSERT(Monitor.MutableColumns.Except(Monitor.InsertedAt, Monitor.UpdatedAt)).
		MODEL(monitor).
		RETURNING(Monitor.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, monitor)
	return monitor, err
}

func (m *MonitorDomain) GetHead(ctx context.Context, organizationID int64) (*model.Monitor, error) {
	stmt := m.listRecursively(organizationID, 1)

	monitor := &model.Monitor{}
	err := stmt.QueryContext(ctx, infra.DB, monitor)
	return monitor, err
}

func (m *MonitorDomain) UpdatePrevious(ctx context.Context, tx *sql.Tx, id, previousID int64) (*model.Monitor, error) {
	now := times.Now()
	monitor := &model.Monitor{
		PrevID:    &previousID,
		UpdatedAt: now,
	}

	updateStmt := Monitor.UPDATE(Monitor.PrevID, Monitor.UpdatedAt).
		MODEL(monitor).WHERE(Monitor.ID.EQ(Int(id))).RETURNING(Monitor.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, monitor)
	return monitor, err
}

func (m *MonitorDomain) listRecursively(organizationID int64, limit int) Statement {
	monitorTree := CTE("monitor_tree")

	stmt := WITH_RECURSIVE(
		monitorTree.AS(
			SELECT(
				Monitor.AllColumns,
			).FROM(
				Monitor,
			).WHERE(
				Monitor.PrevID.IS_NULL(),
			).UNION(
				SELECT(
					Monitor.AllColumns,
				).FROM(
					Monitor.
						INNER_JOIN(monitorTree, Monitor.ID.From(monitorTree).EQ(Monitor.PrevID)),
				),
			),
		),
	)(
		SELECT(
			monitorTree.AllColumns(),
		).FROM(
			monitorTree,
		).WHERE(
			Monitor.OrganizationID.From(monitorTree).EQ(Int(organizationID)),
		).LIMIT(
			int64(limit),
		),
	)

	return stmt
}
