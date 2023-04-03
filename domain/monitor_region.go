package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type MonitorRegionDomain struct{}

func NewMonitorRegionDomain() *MonitorRegionDomain {
	return &MonitorRegionDomain{}
}

func (m *MonitorRegionDomain) GetOldestChecked(
	ctx context.Context,
	monitorID, regionID int64,
) (*model.MonitorRegion, error) {
	stmt := SELECT(MonitorRegion.AllColumns).FROM(MonitorRegion).WHERE(
		MonitorRegion.MonitorID.EQ(Int(monitorID)).AND(MonitorRegion.RegionID.EQ(Int(regionID))),
	).ORDER_BY(MonitorRegion.LastCheckedAt.DESC()).LIMIT(1)

	monitorRegion := &model.MonitorRegion{}
	err := stmt.QueryContext(ctx, infra.DB, monitorRegion)
	return monitorRegion, err
}
