package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

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

func (m *MonitorDomain) list(ctx context.Context, organizationID int64) ([]model.Monitor, error) {
	stmt := m.listRecursively(organizationID)

	var monitors []model.Monitor
	err := stmt.QueryContext(ctx, infra.DB, &monitors)
	return monitors, err
}

func (m *MonitorDomain) getHead(ctx context.Context, organizationID int64) (*model.Monitor, error) {
	stmt := m.listRecursively(organizationID)

	monitor := &model.Monitor{}
	err := stmt.QueryContext(ctx, infra.DB, &monitor)
	return monitor, err
}

func (m *MonitorDomain) listRecursively(organizationID int64) Statement {
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
						INNER_JOIN(monitorTree, Monitor.PrevID.From(monitorTree).EQ(Monitor.ID)),
				),
			),
		),
	)(
		SELECT(
			monitorTree.AllColumns(),
		).FROM(
			monitorTree,
		).WHERE(
			Monitor.OrganizationID.EQ(Int(organizationID)),
		),
	)

	return stmt
}
