package service

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type AlarmPolicyService struct {
	alarmPolicyDomain *domain.AlarmPolicyDomain
}

func NewAlarmPolicyService(alarmPolicyDomain *domain.AlarmPolicyDomain) *AlarmPolicyService {
	return &AlarmPolicyService{alarmPolicyDomain: alarmPolicyDomain}
}

func (a *AlarmPolicyService) Get(
	ctx context.Context, monitorID, organizationID int64,
) (*model.AlarmPolicy, error) {
	alarmPolicy, err := a.alarmPolicyDomain.GetMonitorAlarmPolicy(ctx, organizationID, monitorID)
	if err != nil {
		return a.alarmPolicyDomain.GetGlobal(ctx, organizationID)
	}
	return alarmPolicy, nil
}
