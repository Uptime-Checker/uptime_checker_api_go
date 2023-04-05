package domain

import (
	"context"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type DailyReportDomain struct{}

func NewDailyReportDomain() *DailyReportDomain {
	return &DailyReportDomain{}
}

func (u *DailyReportDomain) Get(ctx context.Context, monitorID int64, date time.Time) (*model.DailyReport, error) {
	stmt := SELECT(DailyReport.AllColumns).FROM(DailyReport).WHERE(
		DailyReport.MonitorID.EQ(Int(monitorID)).AND(DailyReport.Date.EQ(DateT(date))),
	).LIMIT(1)

	dailyReport := &model.DailyReport{}
	err := stmt.QueryContext(ctx, infra.DB, dailyReport)
	return dailyReport, err
}

func (u *DailyReportDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	dailyReport *model.DailyReport,
) (*model.DailyReport, error) {
	insertStmt := DailyReport.INSERT(DailyReport.MutableColumns.Except(DailyReport.InsertedAt, DailyReport.UpdatedAt)).
		MODEL(dailyReport).ON_CONFLICT(DailyReport.Date, DailyReport.MonitorID).DO_NOTHING().
		RETURNING(DailyReport.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, dailyReport)
	return dailyReport, err
}

func (u *DailyReportDomain) Update(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	dailyReport *model.DailyReport,
) (*model.DailyReport, error) {
	updateStmt := DailyReport.
		UPDATE(DailyReport.SuccessfulChecks, DailyReport.ErrorChecks, DailyReport.Downtime, DailyReport.UpdatedAt).
		MODEL(dailyReport).WHERE(DailyReport.ID.EQ(Int(id))).RETURNING(DailyReport.AllColumns)
	err := updateStmt.QueryContext(ctx, tx, dailyReport)
	return dailyReport, err
}
