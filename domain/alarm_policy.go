package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type AlarmPolicyDomain struct{}

func NewAlarmPolicyDomain() *AlarmPolicyDomain {
	return &AlarmPolicyDomain{}
}

func (a *AlarmPolicyDomain) GetGlobal(ctx context.Context, organizationID int64) (*model.AlarmPolicy, error) {
	stmt := SELECT(AlarmPolicy.AllColumns).FROM(AlarmPolicy).
		WHERE(AlarmPolicy.OrganizationID.EQ(Int(organizationID)).AND(AlarmPolicy.MonitorID.IS_NULL())).LIMIT(1)

	policy := &model.AlarmPolicy{}
	err := stmt.QueryContext(ctx, infra.DB, policy)
	return policy, err
}

func (a *AlarmPolicyDomain) GetMonitorAlarmPolicy(
	ctx context.Context,
	organizationID, monitorID int64,
) (*model.AlarmPolicy, error) {
	stmt := SELECT(AlarmPolicy.AllColumns).FROM(AlarmPolicy).
		WHERE(
			AlarmPolicy.OrganizationID.EQ(Int(organizationID)).
				AND(AlarmPolicy.MonitorID.EQ(Int(monitorID))),
		).LIMIT(1)

	policy := &model.AlarmPolicy{}
	err := stmt.QueryContext(ctx, infra.DB, policy)
	return policy, err
}

func (a *AlarmPolicyDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	alarmPolicy *model.AlarmPolicy,
	alarmPolicyName resource.AlarmPolicyName,
) (*model.AlarmPolicy, error) {
	if !alarmPolicyName.Valid() {
		return nil, constant.ErrInvalidAlarmPolicy
	}
	alarmPolicy.Reason = string(alarmPolicyName)

	insertStmt := AlarmPolicy.INSERT(AlarmPolicy.MutableColumns.Except(AlarmPolicy.InsertedAt, AlarmPolicy.UpdatedAt)).
		MODEL(alarmPolicy).RETURNING(AlarmPolicy.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, alarmPolicy)
	return alarmPolicy, err
}
