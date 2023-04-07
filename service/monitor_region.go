package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type MonitorRegionService struct {
	monitorRegionDomain *domain.MonitorRegionDomain
}

func NewMonitorRegionService(monitorRegionDomain *domain.MonitorRegionDomain) *MonitorRegionService {
	return &MonitorRegionService{monitorRegionDomain: monitorRegionDomain}
}

func (m *MonitorRegionService) FirstOrCreate(
	ctx context.Context,
	tx *sql.Tx,
	monitorID, regionID int64,
) (*model.MonitorRegion, error) {
	tracingID := pkg.GetTracingID(ctx)
	monitorRegion, err := m.monitorRegionDomain.GetMonitorRegion(ctx, monitorID, regionID)
	if err == nil {
		lgr.Print(tracingID, 1, "got monitor region, monitor", monitorID, "region", regionID)
		return monitorRegion, nil
	}
	now := times.Now()
	monitorRegion = &model.MonitorRegion{
		Down:          false,
		LastCheckedAt: &now,
		MonitorID:     monitorID,
		RegionID:      regionID,
	}
	lgr.Print(tracingID, 2, "creating monitor region, monitor", monitorID, "region", regionID)
	return m.monitorRegionDomain.Create(ctx, tx, monitorRegion)
}
