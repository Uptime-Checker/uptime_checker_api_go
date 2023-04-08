package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

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
