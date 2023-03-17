package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

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
