package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type AlarmChannelDomain struct{}

func NewAlarmChannelDomain() *AlarmChannelDomain {
	return &AlarmChannelDomain{}
}

func (a *AlarmChannelDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	alarmChannel *model.AlarmChannel,
) (*model.AlarmChannel, error) {
	insertStmt := AlarmChannel.INSERT(AlarmChannel.MutableColumns.
		Except(AlarmChannel.InsertedAt, AlarmChannel.UpdatedAt)).
		MODEL(alarmChannel).RETURNING(AlarmChannel.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, alarmChannel)
	return alarmChannel, err
}

func (a *AlarmChannelDomain) ListByOrganization(
	ctx context.Context,
	organizationID int64,
) ([]model.AlarmChannel, error) {
	stmt := SELECT(AlarmChannel.AllColumns).FROM(AlarmChannel).
		WHERE(AlarmChannel.OrganizationID.EQ(Int(organizationID)))

	var alarmChannels []model.AlarmChannel
	err := stmt.QueryContext(ctx, infra.DB, &alarmChannels)
	return alarmChannels, err
}

func (a *AlarmChannelDomain) ListByMonitor(
	ctx context.Context,
	monitorID int64,
) ([]model.AlarmChannel, error) {
	stmt := SELECT(AlarmChannel.AllColumns).FROM(AlarmChannel).
		WHERE(AlarmChannel.MonitorID.EQ(Int(monitorID)))

	var alarmChannels []model.AlarmChannel
	err := stmt.QueryContext(ctx, infra.DB, &alarmChannels)
	return alarmChannels, err
}
