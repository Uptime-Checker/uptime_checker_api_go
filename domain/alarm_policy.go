package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type AlarmPolicyDomain struct{}

func NewAlarmPolicy() *AlarmPolicyDomain {
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
