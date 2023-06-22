package domain

import (
	"context"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
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

func (m *MonitorDomain) Get(
	ctx context.Context,
	monitorID int64,
) (*model.Monitor, error) {
	stmt := SELECT(Monitor.AllColumns).FROM(Monitor).WHERE(Monitor.ID.EQ(Int(monitorID))).LIMIT(1)

	monitor := &model.Monitor{}
	err := stmt.QueryContext(ctx, infra.DB, monitor)
	return monitor, err
}

func (m *MonitorDomain) GetAll(
	ctx context.Context,
	monitorID int64,
) (*pkg.MonitorWithResources, error) {
	stmt := SELECT(
		Monitor.AllColumns,
		Region.AllColumns,
		Assertion.AllColumns,
	).
		FROM(
			Monitor.
				LEFT_JOIN(MonitorRegion, MonitorRegion.MonitorID.EQ(Monitor.ID)).
				LEFT_JOIN(Region, MonitorRegion.RegionID.EQ(Region.ID)).
				LEFT_JOIN(Assertion, Assertion.MonitorID.EQ(Monitor.ID)),
		).WHERE(Monitor.ID.EQ(Int(monitorID))).LIMIT(1)

	monitor := &pkg.MonitorWithResources{}
	err := stmt.QueryContext(ctx, infra.DB, monitor)
	return monitor, err
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
	bodyFormat resource.MonitorBodyFormat,
) (*model.Monitor, error) {
	if !monitorType.Valid() {
		return nil, constant.ErrInvalidMonitorType
	}
	if !bodyFormat.Valid() {
		return nil, constant.ErrInvalidMonitorBodyFormat
	}

	monitor.Type = int32(monitorType)
	monitor.BodyFormat = int32(bodyFormat)
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

func (m *MonitorDomain) UpdateNext(ctx context.Context, tx *sql.Tx, id, nextID int64) (*model.Monitor, error) {
	now := times.Now()
	monitor := &model.Monitor{
		NextID:    &nextID,
		UpdatedAt: now,
	}

	updateStmt := Monitor.UPDATE(Monitor.NextID, Monitor.UpdatedAt).
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
				Monitor.NextID.IS_NULL(),
			).UNION(
				SELECT(
					Monitor.AllColumns,
				).FROM(
					Monitor.
						INNER_JOIN(monitorTree, Monitor.ID.From(monitorTree).EQ(Monitor.NextID)),
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

func (m *MonitorDomain) ListMonitorsToRun(ctx context.Context, prev, later time.Time) ([]model.Monitor, error) {
	condition := Monitor.NextCheckAt.GT(TimestampT(prev)).
		AND(Monitor.NextCheckAt.LT(TimestampT(later))).
		AND(Monitor.LastCheckedAt.LT(TimestampT(prev)))

	condition = condition.OR(Monitor.LastCheckedAt.IS_NULL())
	condition = condition.AND(Monitor.On.EQ(Bool(true)))

	stmt := SELECT(Monitor.AllColumns).FROM(Monitor).WHERE(condition)

	var monitors []model.Monitor
	err := stmt.QueryContext(ctx, infra.DB, &monitors)
	return monitors, err
}

func (m *MonitorDomain) UpdateOnStatusAndCheckedAt(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	on bool,
	status resource.MonitorStatus,
	lastCheckedAt, nextCheckAt *time.Time,
) (*model.Monitor, error) {
	if !status.Valid() {
		return nil, constant.ErrInvalidMonitorStatus
	}
	now := times.Now()
	monitor := &model.Monitor{
		On:            on,
		Status:        int32(status),
		LastCheckedAt: lastCheckedAt,
		NextCheckAt:   nextCheckAt,
		UpdatedAt:     now,
	}

	updateStmt := Monitor.
		UPDATE(Monitor.On, Monitor.Status, Monitor.LastCheckedAt, Monitor.NextCheckAt, Monitor.UpdatedAt).
		MODEL(monitor).WHERE(Monitor.ID.EQ(Int(id))).RETURNING(Monitor.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, monitor)
	return monitor, err
}

// UpdateNextCheckAt does not need transaction
func (m *MonitorDomain) UpdateNextCheckAt(
	ctx context.Context,
	id int64,
	lastCheckedAt, nextCheckAt *time.Time,
) (*model.Monitor, error) {
	now := times.Now()
	monitor := &model.Monitor{
		LastCheckedAt: lastCheckedAt,
		NextCheckAt:   nextCheckAt,
		UpdatedAt:     now,
	}

	updateStmt := Monitor.UPDATE(Monitor.LastCheckedAt, Monitor.NextCheckAt, Monitor.UpdatedAt).
		MODEL(monitor).WHERE(Monitor.ID.EQ(Int(id))).RETURNING(Monitor.AllColumns)

	err := updateStmt.QueryContext(ctx, infra.DB, monitor)
	return monitor, err
}

func (m *MonitorDomain) UpdateConsecutive(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	status resource.MonitorStatus,
	consecutiveCount int32,
	lastFailedAt *time.Time,
) (*model.Monitor, error) {
	if !status.Valid() {
		return nil, constant.ErrInvalidMonitorStatus
	}
	now := times.Now()
	monitor := &model.Monitor{
		ConsecutiveCount: consecutiveCount,
		Status:           int32(status),
		LastFailedAt:     lastFailedAt,
		UpdatedAt:        now,
	}

	updateStmt := Monitor.
		UPDATE(Monitor.ConsecutiveCount, Monitor.Status, Monitor.LastFailedAt, Monitor.UpdatedAt).
		MODEL(monitor).WHERE(Monitor.ID.EQ(Int(id))).RETURNING(Monitor.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, monitor)
	return monitor, err
}
