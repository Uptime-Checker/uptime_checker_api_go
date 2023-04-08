package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type AlarmDomain struct{}

func NewAlarmDomain() *AlarmDomain {
	return &AlarmDomain{}
}

func (a *AlarmDomain) GetOngoing(ctx context.Context, monitorID int64) (*model.Alarm, error) {
	stmt := SELECT(Alarm.AllColumns).FROM(Alarm).
		WHERE(Alarm.MonitorID.EQ(Int(monitorID)).AND(Alarm.Ongoing.IS_TRUE())).LIMIT(1)

	alarm := &model.Alarm{}
	err := stmt.QueryContext(ctx, infra.DB, alarm)
	return alarm, err
}

func (a *AlarmDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	alarm *model.Alarm,
) (*model.Alarm, error) {
	insertStmt := Alarm.INSERT(Alarm.MutableColumns.Except(Alarm.InsertedAt, Alarm.UpdatedAt)).
		MODEL(alarm).
		RETURNING(Alarm.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, alarm)
	return alarm, err
}

func (a *AlarmDomain) Resolve(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	resolvedByCheckID int64,
) (*model.Alarm, error) {
	now := times.Now()
	alarm := &model.Alarm{
		Ongoing:           false,
		ResolvedByCheckID: &resolvedByCheckID,
		ResolvedAt:        &now,
		UpdatedAt:         now,
	}
	updateStmt := Alarm.UPDATE(
		Alarm.Ongoing, Alarm.ResolvedByCheckID, Alarm.ResolvedAt, Alarm.UpdatedAt,
	).MODEL(alarm).
		WHERE(Alarm.ID.EQ(Int(id))).
		RETURNING(Alarm.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, alarm)
	return alarm, err
}
