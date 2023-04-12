package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type MonitorRegionDomain struct{}

func NewMonitorRegionDomain() *MonitorRegionDomain {
	return &MonitorRegionDomain{}
}

func (m *MonitorRegionDomain) GetOldestChecked(ctx context.Context, monitorID int64) (*model.MonitorRegion, error) {
	stmt := SELECT(MonitorRegion.AllColumns).FROM(MonitorRegion).WHERE(
		MonitorRegion.MonitorID.EQ(Int(monitorID)),
	).ORDER_BY(MonitorRegion.LastCheckedAt.DESC()).LIMIT(1)

	monitorRegion := &model.MonitorRegion{}
	err := stmt.QueryContext(ctx, infra.DB, monitorRegion)
	return monitorRegion, err
}

func (m *MonitorRegionDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	monitorRegion *model.MonitorRegion,
) (*model.MonitorRegion, error) {
	insertStmt := MonitorRegion.INSERT(MonitorRegion.MutableColumns.
		Except(MonitorRegion.InsertedAt, MonitorRegion.UpdatedAt)).
		MODEL(monitorRegion).
		RETURNING(MonitorRegion.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, monitorRegion)
	return monitorRegion, err
}

func (m *MonitorRegionDomain) GetAll(
	ctx context.Context,
	monitorID int64,
) ([]model.MonitorRegion, error) {
	stmt := SELECT(MonitorRegion.AllColumns).FROM(MonitorRegion).WHERE(MonitorRegion.MonitorID.EQ(Int(monitorID)))

	var monitorRegions []model.MonitorRegion
	err := stmt.QueryContext(ctx, infra.DB, &monitorRegions)
	return monitorRegions, err
}

func (m *MonitorRegionDomain) GetMonitorRegion(
	ctx context.Context,
	monitorID, regionID int64,
) (*model.MonitorRegion, error) {
	stmt := SELECT(MonitorRegion.AllColumns).FROM(MonitorRegion).WHERE(
		MonitorRegion.MonitorID.EQ(Int(monitorID)).AND(MonitorRegion.RegionID.EQ(Int(regionID))),
	).LIMIT(1)

	monitorRegion := &model.MonitorRegion{}
	err := stmt.QueryContext(ctx, infra.DB, monitorRegion)
	return monitorRegion, err
}

func (m *MonitorRegionDomain) GetWithAllAssoc(
	ctx context.Context,
	monitorRegionID int64,
) (*pkg.MonitorRegionWithAssertions, error) {
	stmt := SELECT(
		MonitorRegion.AllColumns,
		Monitor.AllColumns,
		Region.AllColumns,
		Assertion.AllColumns,
	).
		FROM(
			MonitorRegion.
				LEFT_JOIN(Monitor, MonitorRegion.MonitorID.EQ(Monitor.ID)).
				LEFT_JOIN(Region, MonitorRegion.RegionID.EQ(Region.ID)).
				LEFT_JOIN(Assertion, MonitorRegion.MonitorID.EQ(Assertion.MonitorID)),
		).
		WHERE(MonitorRegion.ID.EQ(Int(monitorRegionID)))

	monitorRegionWithAssertions := &pkg.MonitorRegionWithAssertions{}
	err := stmt.QueryContext(ctx, infra.DB, monitorRegionWithAssertions)
	return monitorRegionWithAssertions, err
}

func (m *MonitorRegionDomain) UpdateDown(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	down bool,
) (*model.MonitorRegion, error) {
	now := times.Now()
	monitorRegion := &model.MonitorRegion{
		Down:          down,
		LastCheckedAt: &now,
		UpdatedAt:     now,
	}

	updateStmt := MonitorRegion.UPDATE(MonitorRegion.Down, MonitorRegion.LastCheckedAt, Monitor.UpdatedAt).
		MODEL(monitorRegion).WHERE(MonitorRegion.ID.EQ(Int(id))).RETURNING(MonitorRegion.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, monitorRegion)
	return monitorRegion, err
}
