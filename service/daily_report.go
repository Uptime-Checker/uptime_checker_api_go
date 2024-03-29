package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type DailyReportService struct {
	dailyReportDomain *domain.DailyReportDomain
}

func NewDailyReportService(dailyReportDomain *domain.DailyReportDomain) *DailyReportService {
	return &DailyReportService{dailyReportDomain: dailyReportDomain}
}

func (d *DailyReportService) Add(
	ctx context.Context,
	tx *sql.Tx, monitorID, organizationID int64, success bool,
) (*model.DailyReport, error) {
	now := time.Now()
	successfulChecks := 0
	errorChecks := 0

	dailyReport, err := d.dailyReportDomain.Get(ctx, monitorID, now)
	if err != nil {
		// Create
		if success {
			successfulChecks = 1
		} else {
			errorChecks = 1
		}
		dailyReport = &model.DailyReport{
			SuccessfulChecks: int32(successfulChecks),
			ErrorChecks:      int32(errorChecks),
			Downtime:         0, // fresh, huh
			Date:             now,
			MonitorID:        monitorID,
			OrganizationID:   organizationID,
		}
		return d.dailyReportDomain.Create(ctx, tx, dailyReport)
	}
	// Update
	if success {
		dailyReport.SuccessfulChecks++
	} else {
		dailyReport.ErrorChecks--
	}
	dailyReport.UpdatedAt = now
	return d.dailyReportDomain.Update(ctx, tx, dailyReport.ID, dailyReport)
}

func (d *DailyReportService) UpdateDailyDowntime(
	ctx context.Context, tx *sql.Tx, dailyReport *model.DailyReport, from, to time.Time,
) (*model.DailyReport, error) {
	duration := to.Sub(from).Seconds()
	dailyReport.Downtime += int32(duration)
	dailyReport.UpdatedAt = time.Now()
	return d.dailyReportDomain.Update(ctx, tx, dailyReport.ID, dailyReport)
}
