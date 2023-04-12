package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type OrganizationService struct {
	organizationDomain *domain.OrganizationDomain
	alarmPolicyDomain  *domain.AlarmPolicyDomain
}

func NewOrganizationService(
	organizationDomain *domain.OrganizationDomain,
	alarmPolicyDomain *domain.AlarmPolicyDomain,
) *OrganizationService {
	return &OrganizationService{organizationDomain: organizationDomain, alarmPolicyDomain: alarmPolicyDomain}
}

func (o *OrganizationService) CreateOrganizationAlarmPolicy(
	ctx context.Context, tx *sql.Tx, organizationID int64,
) (*model.AlarmPolicy, error) {
	threshold := int32(constant.DefaultOrganizationAlarmErrorThreshold)
	alarmPolicy := &model.AlarmPolicy{
		Threshold:      threshold,
		OrganizationID: organizationID,
	}

	return o.alarmPolicyDomain.Create(ctx, tx, alarmPolicy, resource.AlarmPolicyErrorThreshold)
}
